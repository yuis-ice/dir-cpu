import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'dir-cpu',
  description: 'Real-time CPU usage aggregated by filesystem directory',
  base: '/dir-cpu/',

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/dir-cpu/favicon.svg' }],
  ],

  themeConfig: {
    logo: { svg: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M3 3h18v2H3V3zm0 4h18v2H3V7zm0 4h12v2H3v-2zm0 4h12v2H3v-2zm0 4h18v2H3v-2z"/></svg>' },

    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Architecture', link: '/architecture/overview' },
      { text: 'Recipes', link: '/recipes/filter-by-subtree' },
      {
        text: 'v0.1.0',
        items: [
          { text: 'GitHub', link: 'https://github.com/yuis-ice/dir-cpu' },
          { text: 'pkg.go.dev', link: 'https://pkg.go.dev/github.com/yuis-ice/dir-cpu' },
        ]
      }
    ],

    sidebar: [
      {
        text: 'Guide',
        items: [
          { text: 'Getting Started', link: '/guide/getting-started' },
          { text: 'CLI Reference', link: '/guide/cli-reference' },
          { text: 'cwd vs exe mode', link: '/guide/cwd-vs-exe' },
          { text: 'Permissions', link: '/guide/permissions' },
        ]
      },
      {
        text: 'Architecture',
        items: [
          { text: 'Overview', link: '/architecture/overview' },
          { text: 'CPU Calculation', link: '/architecture/cpu-calculation' },
          { text: 'Tree Roll-up Algorithm', link: '/architecture/tree-rollup' },
          { text: 'Performance & Overhead', link: '/architecture/performance' },
        ]
      },
      {
        text: 'Recipes',
        items: [
          { text: 'Filter by Subtree', link: '/recipes/filter-by-subtree' },
          { text: 'CI/CD Anomaly Detection', link: '/recipes/cicd-anomaly' },
          { text: 'Security & Forensics', link: '/recipes/security-forensics' },
        ]
      },
      {
        text: 'Reference',
        items: [
          { text: 'Comparison: top / htop / glances', link: '/reference/comparison' },
          { text: 'FAQ', link: '/reference/faq' },
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/yuis-ice/dir-cpu' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 yuis-ice'
    },

    editLink: {
      pattern: 'https://github.com/yuis-ice/dir-cpu/edit/main/docs/:path',
      text: 'Edit this page on GitHub'
    },

    search: {
      provider: 'local'
    }
  }
})
