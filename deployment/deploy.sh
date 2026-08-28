#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Check if script is run as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (or use sudo)"
  exit 1
fi

echo "========================================"
echo " Starting Server Deployment Script "
echo "========================================"

# Variables
DB_NAME="invoice_app"
DB_USER="invoice_user"
# Generate a random password for the database
DB_PASSWORD=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 16 | head -n 1)
APP_DIR="/opt/invoice-builder"

# 1. System Update & Dependencies
echo "[1/6] Updating system and installing dependencies (nginx, postgresql, ufw)..."
apt-get update
apt-get install -y nginx postgresql postgresql-contrib ufw

# 2. Firewall Setup (UFW)
echo "[2/6] Configuring Firewall (UFW)..."
ufw allow OpenSSH
ufw allow 'Nginx Full'
ufw --force enable

# 3. PostgreSQL Database Setup
echo "[3/6] Configuring PostgreSQL..."
# Check if database exists, create if not
sudo -u postgres psql -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || sudo -u postgres psql -c "CREATE DATABASE $DB_NAME;"
sudo -u postgres psql -tc "SELECT 1 FROM pg_roles WHERE rolname = '$DB_USER'" | grep -q 1 || sudo -u postgres psql -c "CREATE USER $DB_USER WITH ENCRYPTED PASSWORD '$DB_PASSWORD';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"

# 4. App Directory Setup
echo "[4/6] Setting up application directory..."
mkdir -p $APP_DIR
# Assuming the script is run from the directory containing the binary and .env
# If not, you will need to manually copy these.
if [ -f "invoice-builder" ]; then
    cp invoice-builder $APP_DIR/
    chmod +x $APP_DIR/invoice-builder
else
    echo "Warning: 'invoice-builder' binary not found in current directory. Please copy it to $APP_DIR manually."
fi

echo "Building Vue frontend..."
if [ -d "frontend" ]; then
    cd frontend
    npm install
    npm run build
    cd ..
    mkdir -p $APP_DIR/frontend/dist
    cp -r frontend/dist/* $APP_DIR/frontend/dist/
else
    echo "Error: 'frontend' directory not found."
    exit 1
fi

echo "Setting up uploads directory..."
mkdir -p $APP_DIR/ui/asset/img/uploads
chmod 777 $APP_DIR/ui/asset/img/uploads

if [ -f ".env" ]; then
    cp .env $APP_DIR/
else
    echo "Warning: '.env' file not found in current directory. Creating a template..."
    cat > $APP_DIR/.env <<EOL
APP_PORT=4000
APP_ENV=production
DB_HOST=localhost
DB_PORT=5432
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
DB_SSLMODE=disable
JWT_SECRET=replace_me_with_a_secure_random_string
SMTP_HOST=your_smtp_host
SMTP_PORT=587
SMTP_USERNAME=your_email
SMTP_PASSWORD=your_password
SMTP_SENDER="Teks-Invoice <your_email>"
PAYSTACK_SECRET_KEY=your_paystack_secret
PAYSTACK_PUBLIC_KEY=your_paystack_public
EOL
    echo "Created template .env at $APP_DIR/.env. PLEASE EDIT THIS FILE."
fi

# 5. Systemd Service Setup
echo "[5/6] Configuring Systemd Service..."
echo "Building Go API application..."
go build -o invoice-builder ./cmd/api
if [ -f "invoice-builder.service" ]; then
    cp invoice-builder.service /etc/systemd/system/
    systemctl daemon-reload
    systemctl enable invoice-builder
    systemctl restart invoice-builder
else
    echo "Warning: 'invoice-builder.service' not found. Please place it in /etc/systemd/system/ manually."
fi

# 6. Nginx Setup
echo "[6/6] Configuring Nginx Reverse Proxy..."
if [ -f "nginx.conf" ]; then
    cp nginx.conf /etc/nginx/sites-available/teks-invoice
    # Create symlink if it doesn't exist
    if [ ! -L /etc/nginx/sites-enabled/teks-invoice ]; then
        ln -s /etc/nginx/sites-available/teks-invoice /etc/nginx/sites-enabled/
    fi
    # Remove default nginx site
    rm -f /etc/nginx/sites-enabled/default
    
    # We won't restart Nginx just yet because the Cloudflare certs are missing.
    echo "Nginx configuration copied."
else
    echo "Warning: 'nginx.conf' not found."
fi

echo "========================================"
echo " Deployment Script Finished! "
echo "========================================"
echo ""
echo "!!! ACTION REQUIRED FOR CLOUDFLARE FULL STRICT SSL !!!"
echo "1. Go to your Cloudflare Dashboard -> SSL/TLS -> Origin Server."
echo "2. Click 'Create Certificate'."
echo ""
echo "3. We will now open an editor for the ORIGIN CERTIFICATE."
echo "   -> Paste the certificate contents."
echo "   -> Press Ctrl+X, then type Y, then press Enter to save."
read -p "Press Enter to open the editor..."
nano /etc/ssl/certs/cloudflare-origin.pem

echo ""
echo "4. We will now open an editor for the PRIVATE KEY."
echo "   -> Paste the private key contents."
echo "   -> Press Ctrl+X, then type Y, then press Enter to save."
read -p "Press Enter to open the editor..."
nano /etc/ssl/private/cloudflare-origin.key

echo ""
echo "5. Restarting Nginx to apply certificates..."
systemctl restart nginx
echo ""
echo "Your generated Database Password is: $DB_PASSWORD"
echo "(It has been automatically saved in $APP_DIR/.env if the script generated the file)"
