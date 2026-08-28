#!/bin/bash

# ============================================================
# Teardown Script for Invoice Builder (Server-Side)
# Run this on the server as root: sudo bash teardown.sh
# ============================================================

set -e

# Safety check
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (or use sudo)"
  exit 1
fi

echo "========================================"
echo " Invoice Builder - Full Teardown Script "
echo "========================================"
echo ""
echo "This will PERMANENTLY remove:"
echo "  - The invoice-builder systemd service"
echo "  - The application binary and .env at /opt/invoice-builder/"
echo "  - The Nginx site configuration"
echo "  - The Cloudflare Origin SSL certificate and private key"
echo "  - The staging copy at ~/invoice-setup/"
echo ""
read -p "Are you SURE you want to proceed? (type 'yes' to confirm): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
  echo "Aborted."
  exit 0
fi

# 1. Stop and disable the systemd service
echo ""
echo "[1/5] Stopping and removing systemd service..."
if systemctl is-active --quiet invoice-builder 2>/dev/null; then
  systemctl stop invoice-builder
  echo "  -> Service stopped."
fi
if systemctl is-enabled --quiet invoice-builder 2>/dev/null; then
  systemctl disable invoice-builder
  echo "  -> Service disabled."
fi
if [ -f /etc/systemd/system/invoice-builder.service ]; then
  rm -f /etc/systemd/system/invoice-builder.service
  systemctl daemon-reload
  echo "  -> Service file removed and daemon reloaded."
else
  echo "  -> Service file not found, skipping."
fi

# 2. Remove the application directory
echo ""
echo "[2/5] Removing application directory /opt/invoice-builder/..."
if [ -d /opt/invoice-builder ]; then
  rm -rf /opt/invoice-builder
  echo "  -> Removed /opt/invoice-builder/"
else
  echo "  -> Directory not found, skipping."
fi

# Also remove the staging copy if it exists
STAGING_DIR="/home/bryt/invoice-setup"
if [ -d "$STAGING_DIR" ]; then
  rm -rf "$STAGING_DIR"
  echo "  -> Removed staging directory $STAGING_DIR"
fi

# 3. Remove Nginx configuration
echo ""
echo "[3/5] Removing Nginx configuration..."
if [ -L /etc/nginx/sites-enabled/teks-invoice ]; then
  rm -f /etc/nginx/sites-enabled/teks-invoice
  echo "  -> Removed sites-enabled symlink."
fi
if [ -f /etc/nginx/sites-available/teks-invoice ]; then
  rm -f /etc/nginx/sites-available/teks-invoice
  echo "  -> Removed sites-available config."
fi
# Restore default nginx site so nginx doesn't fail on restart
if [ -f /etc/nginx/sites-available/default ] && [ ! -L /etc/nginx/sites-enabled/default ]; then
  ln -s /etc/nginx/sites-available/default /etc/nginx/sites-enabled/default 2>/dev/null || true
  echo "  -> Restored default Nginx site."
fi
# Test and reload nginx
if command -v nginx &>/dev/null; then
  nginx -t 2>/dev/null && systemctl reload nginx
  echo "  -> Nginx reloaded."
fi

# 4. Delete Cloudflare Origin SSL certificate and private key
echo ""
echo "[4/5] Removing Cloudflare Origin SSL certificate and key..."
if [ -f /etc/ssl/certs/cloudflare-origin.pem ]; then
  # Securely wipe the certificate
  shred -u /etc/ssl/certs/cloudflare-origin.pem 2>/dev/null || rm -f /etc/ssl/certs/cloudflare-origin.pem
  echo "  -> Removed /etc/ssl/certs/cloudflare-origin.pem"
else
  echo "  -> Certificate not found, skipping."
fi
if [ -f /etc/ssl/private/cloudflare-origin.key ]; then
  # Securely wipe the private key
  shred -u /etc/ssl/private/cloudflare-origin.key 2>/dev/null || rm -f /etc/ssl/private/cloudflare-origin.key
  echo "  -> Removed /etc/ssl/private/cloudflare-origin.key"
else
  echo "  -> Private key not found, skipping."
fi

# 5. Clean up syslog entries (optional)
echo ""
echo "[5/5] Cleaning up logs..."
if command -v journalctl &>/dev/null; then
  journalctl --rotate 2>/dev/null || true
  journalctl --vacuum-time=1s --unit=invoice-builder 2>/dev/null || true
  echo "  -> Rotated and vacuumed invoice-builder journal logs."
fi

echo ""
echo "========================================"
echo " Teardown Complete!                     "
echo "========================================"
echo ""
echo "NOTE: The PostgreSQL database 'invoice_app' was NOT deleted."
echo "If you want to drop it too, run:"
echo "  sudo -u postgres psql -c \"DROP DATABASE invoice_app;\""
echo "  sudo -u postgres psql -c \"DROP USER invoice_user;\""
echo ""
echo "Also remember to revoke the Cloudflare Origin Certificate"
echo "from the Cloudflare Dashboard -> SSL/TLS -> Origin Server."
