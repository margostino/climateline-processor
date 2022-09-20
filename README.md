# climateline-processor

[![Build & Test](https://github.com/margostino/climateline-processor/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/margostino/climateline-processor/actions/workflows/main.yml)
[![fetcher-job](https://github.com/margostino/climateline-processor/actions/workflows/job.yml/badge.svg?branch=master)](https://github.com/margostino/climateline-processor/actions/workflows/job.yml)

This repo implements the Job to fetch the news of the day related to Climate Change.  
The Job retrieve every news into the Admin Telegram Bot in order to approve, edit and push new article into [Climateline](https://climateline.vercel.app/) website.

### Architecture

![](assets/architecture.png#100x)

### Features

- [x] Cron job for fetching news
- [x] Telegram bot for handling news workflow (upload, update) towards [Climateline](https://climateline.vercel.app/)
- [x] Basic and short-time caching
- [ ] Multi-Source News Fetcher
- [ ] Create entirely new articles
- [ ] Detect duplicated article before pushing
- [ ] Pushing by replying

