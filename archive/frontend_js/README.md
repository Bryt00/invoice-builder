# Vue 3 + Vite + Pinia + Tailwind CSS + daisyUI + Lucide Icons

This template provides a solid foundation for building scalable Vue 3 applications with modern tooling:

- **Vue 3 + Vite**: Blazing fast development and build tooling
- **Pinia**: Official state management for Vue
- **Tailwind CSS**: Utility-first CSS framework
- **daisyUI**: UI component library for Tailwind CSS
- **Lucide Icons**: Beautiful, open-source icon set

## Features

- ✅ Ready to use with Pinia stores for auth, profile, invoices, clients, etc.
- ✅ Pre-configured Tailwind CSS with daisyUI components
- ✅ Lucide Icons available throughout the app
- ✅ Project structure following Vue ecosystem best practices
- ✅ Responsive design out of the box

## Getting Started

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start development server:
   ```bash
   npm run dev
   ```
   The app will be available at `http://localhost:5173`

3. Run production build:
   ```bash
   npm run build
   ```

## Available Scripts

| Script | Description |
|--------|-------------|
| `dev` | Start development server |
| `build` | Build for production |
| `preview` | Preview production build |
| `lint` | Run ESLint to check for code issues |

## Project Structure

```
frontend/
├── src/
│   ├── components/       # Reusable UI components
│   ├── layouts/          # App layout components
│   ├── pages/            # Page components organized by feature
│   │   ├── auth/         # Authentication pages
│   │   ├── app/          # Main application pages
│   │   ├── public/       # Public pages
│   │   └── admin/        # Admin pages
│   ├── stores/           # Pinia stores
│   ├── utils/            # Utility functions and helpers
│   ├── App.vue           # Root component
│   └── main.js           # Entry point
├── public/               # Static assets
├── index.html            # HTML entry point
├── package.json
├── vite.config.js
└── tailwind.config.js
```

## Tailwind CSS & daisyUI

Tailwind CSS is already configured with daisyUI for easy access to pre-built UI components. You can find all available components in the [daisyUI documentation](https://daisyui.com/components).

**Usage Example:**

```vue
<button class="btn btn-primary">
  Primary Button
</button>

<div class="alert alert-success">
  Success message
</div>

<div class="card w-96 bg-base-100 shadow-xl">
  <div class="card-body">
    <h2 class="card-title">Card Title</h2>
    <p>Some card content here...</p>
    <div class="card-actions justify-end">
      <button class="btn btn-primary">Buy Now</button>
    </div>
  </div>
</div>
```

## Lucide Icons

Lucide icons are available throughout the application. Use the `<Icon>` component like this:

```vue
<script setup>
import Icon from '@/components/Icon.vue'
</script>

<template>
  <div>
    <Icon name="home" class="w-6 h-6" />
    <Icon name="user" class="w-5 h-5 text-blue-500" />
    <Icon name="settings" class="w-8 h-8 text-green-600" />
  </div>
</template>
```

**Available Icon Names:**
[Lucide Icons](https://lucide.dev/icons) (use the kebab-case name without the `lucide-` prefix)

## API Integration

The app is configured with a proxy to forward API requests to your backend:

```js
// vite.config.js
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:4000',
      changeOrigin: true
    }
  }
}
```

**Usage in components:**

```js
const response = await fetch('/api/invoices')
```

## Routing

The router is set up with Vue Router v4 and includes:
- Authentication guards (`requiresAuth`, `guestOnly`)
- Profile completion checks (`requiresProfile`)
- Admin role checks (`requiresAdmin`)
- Organized routes by feature

## Pinia Stores

All stores are centralized in the `src/stores/` directory:

```javascript
// src/stores/auth.js
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    accessToken: null,
    isAuthenticated: false
  }),
  // ...
})

// src/stores/profile.js
export const useProfileStore = defineStore('profile', {
  state: () => ({
    profile: null
  }),
  // ...
})

// src/stores/invoice.js
export const useInvoiceStore = defineStore('invoice', {
  state: () => ({
    invoices: [],
    loading: false
  }),
  // ...
})

// src/stores/client.js
export const useClientStore = defineStore('client', {
  state: () => ({
    clients: [],
    loading: false
  }),
  // ...
})
```

## Development Setup

1. **Install dependencies**:
   ```bash
   npm install
   ```

2. **Start dev server**:
   ```bash
   npm run dev
   ```

3. **View all available commands**:
   ```bash
   npm run
   ```

4. **Run tests** (if available):
   ```bash
   npm test
   ```

## Scaling Up

When your application grows, you can maintain organization by:
- Adding new feature modules to the `pages/` directory
- Creating new Pinia stores in `src/stores/`
- Organizing reusable components in `src/components/`
- Using composition functions in `src/composables/`

## License

This template is provided under the MIT License. Feel free to use, modify, and distribute it for your projects.
