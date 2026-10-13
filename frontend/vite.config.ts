import { fileURLToPath, URL } from 'node:url'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import { VitePWA } from 'vite-plugin-pwa'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    VitePWA({
      registerType: 'autoUpdate',
      injectRegister: null,
      includeAssets: [
        'favicon.png',
        'pwa-192x192.png',
        'pwa-512x512.png',
        'pwa-maskable-192x192.png',
        'pwa-maskable-512x512.png',
        'apple-touch-icon-180x180.png',
        'icons.svg'
      ],
      manifest: {
        id: '/',
        name: 'Teks-Invoice — Billing for Modern Professionals',
        short_name: 'Teks-Invoice',
        description: 'Fast, precise invoice generator and billing management for freelancers and small businesses.',
        theme_color: '#4a7c59',
        background_color: '#FAF6F0',
        display: 'standalone',
        orientation: 'portrait-primary',
        scope: '/',
        start_url: '/',
        icons: [
          {
            src: '/pwa-192x192.png',
            sizes: '192x192',
            type: 'image/png',
            purpose: 'any'
          },
          {
            src: '/pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any'
          },
          {
            src: '/pwa-maskable-192x192.png',
            sizes: '192x192',
            type: 'image/png',
            purpose: 'maskable'
          },
          {
            src: '/pwa-maskable-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'maskable'
          }
        ],
        shortcuts: [
          {
            name: 'New Invoice',
            short_name: 'New Invoice',
            description: 'Create and issue a new invoice',
            url: '/user/invoices/new',
            icons: [{ src: '/pwa-192x192.png', sizes: '192x192' }]
          },
          {
            name: 'Dashboard',
            short_name: 'Dashboard',
            description: 'View dashboard overview and analytics',
            url: '/user/dashboard',
            icons: [{ src: '/pwa-192x192.png', sizes: '192x192' }]
          },
          {
            name: 'Clients',
            short_name: 'Clients',
            description: 'Manage client directory',
            url: '/user/clients',
            icons: [{ src: '/pwa-192x192.png', sizes: '192x192' }]
          }
        ],
        categories: ['business', 'finance', 'productivity']
      },
      workbox: {
        globPatterns: ['**/*.{js,css,html,ico,png,svg,woff,woff2}'],
        navigateFallback: '/index.html',
        navigateFallbackDenylist: [/^\/api/],
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/fonts\.googleapis\.com\/.*/i,
            handler: 'StaleWhileRevalidate',
            options: {
              cacheName: 'google-fonts-stylesheets'
            }
          },
          {
            urlPattern: /^https:\/\/fonts\.gstatic\.com\/.*/i,
            handler: 'CacheFirst',
            options: {
              cacheName: 'google-fonts-webfonts',
              expiration: {
                maxEntries: 30,
                maxAgeSeconds: 60 * 60 * 24 * 365 // 1 year
              },
              cacheableResponse: {
                statuses: [0, 200]
              }
            }
          },
          {
            urlPattern: /^\/asset\/.*/i,
            handler: 'StaleWhileRevalidate',
            options: {
              cacheName: 'app-vendor-assets',
              cacheableResponse: {
                statuses: [0, 200]
              }
            }
          }
        ]
      }
    })
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    port: 5174,
    proxy: {
      '/api': {
        target: 'http://localhost:4000',
        changeOrigin: true
      },
      '/asset': {
        target: 'http://localhost:4000',
        changeOrigin: true
      }
    }
  }
})

