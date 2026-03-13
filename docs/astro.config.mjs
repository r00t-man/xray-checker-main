// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";

// https://astro.build/config
export default defineConfig({
  site: "https://xray-checker.kutovoy.dev",
  integrations: [
    starlight({
      title: "Xray Checker",
      favicon: "/favicon.svg",
      logo: {
        light: "./src/assets/logo-light.svg",
        dark: "./src/assets/logo-dark.svg",
      },
      head: [
        {
          tag: "link",
          attrs: {
            rel: "icon",
            type: "image/png",
            href: "/favicon-96x96.png",
            sizes: "96x96",
          },
        },
        {
          tag: "link",
          attrs: {
            rel: "shortcut icon",
            href: "/favicon.ico",
          },
        },
        {
          tag: "link",
          attrs: {
            rel: "apple-touch-icon",
            sizes: "180x180",
            href: "/apple-touch-icon.png",
          },
        },
        {
          tag: "meta",
          attrs: {
            name: "apple-mobile-web-app-title",
            content: "Xray Checker",
          },
        },
        {
          tag: "link",
          attrs: {
            rel: "manifest",
            href: "/site.webmanifest",
          },
        },
        // Plausible Analytics
        {
          tag: "script",
          attrs: {
            async: true,
            src: "https://ps.log.rw/js/pa-mlJHNSq4iSgTf0o8D8qJM.js",
          },
        },
        {
          tag: "script",
          content:
            'window.plausible=window.plausible||function(){(plausible.q=plausible.q||[]).push(arguments)},plausible.init=plausible.init||function(i){plausible.o=i||{}};plausible.init()',
        },
      ],
      editLink: {
        baseUrl: "https://github.com/kutovoys/xray-checker/edit/main/docs/",
      },
      customCss: ["./src/styles/custom.css"],
      social: {
        github: "https://github.com/kutovoys/xray-checker",
        telegram: "https://t.me/+VEzFQmaTZcQ5ZGYy",
        linkedin: "https://www.linkedin.com/in/kutovoys/",
      },
      defaultLocale: "root",
      locales: {
        root: {
          label: "English",
          lang: "en",
        },
        ru: {
          label: "Русский",
          lang: "ru",
        },
        fa: {
          label: "فارسی",
          lang: "fa",
          dir: "rtl",
        },
      },
      sidebar: [
        {
          label: "Introduction",
          translations: {
            ru: "Введение",
            fa: "مقدمه",
          },
          items: [
            {
              label: "Overview",
              translations: {
                ru: "Обзор",
                fa: "معرفی",
              },
              slug: "index",
            },
            {
              label: "Features",
              translations: {
                ru: "Возможности",
                fa: "امکانات",
              },
              slug: "intro/features",
            },
            {
              label: "Architecture",
              translations: {
                ru: "Архитектура",
                fa: "معماری",
              },
              slug: "intro/architecture",
            },
            {
              label: "Quick Start",
              translations: {
                ru: "Быстрый старт",
                fa: "شروع سریع",
              },
              slug: "intro/quick-start",
            },
          ],
        },
        {
          label: "Usage",
          translations: {
            ru: "Использование",
            fa: "استفاده",
          },
          items: [
            {
              label: "CLI",
              translations: {
                ru: "CLI",
                fa: "خط فرمان",
              },
              slug: "usage/cli",
            },
            {
              label: "Docker",
              translations: {
                ru: "Docker",
                fa: "داکر",
              },
              slug: "usage/docker",
            },
            {
              label: "GitHub Actions",
              translations: {
                ru: "GitHub Actions",
                fa: "GitHub Actions",
              },
              slug: "usage/github-actions",
            },
            {
              label: "API Reference",
              translations: {
                ru: "API Reference",
                fa: "مرجع API",
              },
              slug: "usage/api-reference",
            },
            {
              label: "Troubleshooting",
              translations: {
                ru: "Устранение неполадок",
                fa: "عیب‌یابی",
              },
              slug: "usage/troubleshooting",
            },
          ],
        },
        {
          label: "Configuration",
          translations: {
            ru: "Конфигурация",
            fa: "پیکربندی",
          },
          items: [
            {
              label: "Environment Variables",
              translations: {
                ru: "Переменные окружения",
                fa: "متغیرهای محیطی",
              },
              slug: "configuration/envs",
            },
            {
              label: "Subscription Format",
              translations: {
                ru: "Формат подписки",
                fa: "فرمت اشتراک",
              },
              slug: "configuration/subscription",
            },
            {
              label: "Check Methods",
              translations: {
                ru: "Методы проверки",
                fa: "روش‌های بررسی",
              },
              slug: "configuration/check-methods",
            },
            {
              label: "Advanced Configuration",
              translations: {
                ru: "Расширенная конфигурация",
                fa: "پیکربندی پیشرفته",
              },
              slug: "configuration/advanced-conf",
            },
            {
              label: "Public Status Page",
              translations: {
                ru: "Публичная страница статуса",
                fa: "صفحه وضعیت عمومی",
              },
              slug: "configuration/status-page",
            },
            {
              label: "Web Customization",
              translations: {
                ru: "Кастомизация веб-интерфейса",
                fa: "سفارشی‌سازی وب",
              },
              slug: "configuration/web-customization",
              badge: { text: "NEW", variant: "success" },
            },
          ],
        },
        {
          label: "Integrations",
          translations: {
            ru: "Интеграции",
            fa: "یکپارچه‌سازی‌ها",
          },
          items: [
            {
              label: "Metrics",
              translations: {
                ru: "Метрики",
                fa: "متریک‌ها",
              },
              slug: "integrations/metrics",
            },
            {
              label: "Prometheus Setup",
              translations: {
                ru: "Настройка Prometheus",
                fa: "راه‌اندازی Prometheus",
              },
              slug: "integrations/prometheus",
            },
            {
              label: "Uptime Kuma",
              translations: {
                ru: "Uptime Kuma",
                fa: "Uptime Kuma",
              },
              slug: "integrations/uptime-kuma",
            },
            {
              label: "Grafana Dashboards",
              translations: {
                ru: "Grafana Dashboard",
                fa: "داشبوردهای Grafana",
              },
              slug: "integrations/grafana",
              badge: { text: "WIP", variant: "caution" },
            },
            {
              label: "Alternatives",
              translations: {
                ru: "Альтернативы",
                fa: "جایگزین‌ها",
              },
              slug: "integrations/alternatives",
            },
          ],
        },
        {
          label: "Badges",
          translations: {
            ru: "Бейджи",
            fa: "نشان‌ها",
          },
          badge: { text: "NEW", variant: "success" },
          items: [
            {
              label: "Overview",
              translations: {
                ru: "Обзор",
                fa: "معرفی",
              },
              slug: "badges/overview",
            },
            {
              label: "Playground",
              translations: {
                ru: "Конструктор",
                fa: "سازنده",
              },
              slug: "badges/playground",
            },
          ],
        },
        {
          label: "Contributing",
          translations: {
            ru: "Участие в разработке",
            fa: "مشارکت",
          },
          items: [
            {
              label: "Development Guide",
              translations: {
                ru: "Руководство для разработчиков",
                fa: "راهنمای توسعه",
              },
              link: "/contributing/guide",
            },
          ],
        },
        {
          label: "Other Software",
          translations: {
            fa: "نرم‌افزارهای دیگر",
          },
          items: [
            {
              label: "Xray Torrent Blocker",
              link: "https://github.com/kutovoys/xray-torrent-blocker",
            },
            {
              label: "Speedtest Exporter",
              link: "https://github.com/kutovoys/speedtest-exporter",
            },
            {
              label: "Marzban Exporter",
              link: "https://github.com/kutovoys/marzban-exporter",
            },
          ],
        },
        {
          label: "VPN Recommendation",
          translations: {
            ru: "Рекомендуем",
            fa: "VPN پیشنهادی",
          },
          items: [
            {
              label: "BlancVPN",
              link: "https://getblancvpn.com/pricing?promo=klugscl&ref=xc-doc",
              badge: { text: "KLUGSCL", variant: "success" },
              attrs: {
                class: "blanc-vpn-menu-item",
              },
            },
          ],
        },
      ],
    }),
  ],
});
