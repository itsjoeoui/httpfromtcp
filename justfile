setup:
  mkdir assets
  curl -o assets/vim.mp4 https://storage.googleapis.com/qvault-webapp-dynamic-assets/lesson_videos/vim-vs-neovim-prime.mp4

run:
  go run ./cmd/httpserver
