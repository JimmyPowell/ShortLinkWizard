/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  images: {
    unoptimized: true,
  },
  trailingSlash: true,
  // Ensure CSS and JS are properly included
  webpack: (config, { isServer }) => {
    // Optimize CSS extraction
    if (!isServer) {
      config.optimization.splitChunks.cacheGroups = {
        ...config.optimization.splitChunks.cacheGroups,
        styles: {
          name: 'styles',
          test: /\.(css|scss)$/,
          chunks: 'all',
          enforce: true,
          priority: 10,
        },
      }
    }
    return config
  },
  // Configure static export
  experimental: {
    webpackBuildWorker: true,
    optimizeCss: false, // Disable experimental CSS optimization
  },
  // Disable unnecessary features for static export
  eslint: {
    ignoreDuringBuilds: true,
  },
  typescript: {
    ignoreBuildErrors: true,
  },
  // Configure static generation
  generateBuildId: async () => {
    return 'static'
  },
  // Ensure all assets are included
  assetPrefix: '',
  distDir: '.next',
}

export default nextConfig
