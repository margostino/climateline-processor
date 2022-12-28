# climateline-processor

[![Build & Test](https://github.com/margostino/climateline-processor/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/margostino/climateline-processor/actions/workflows/main.yml)
[![collector-job](https://github.com/margostino/climateline-processor/actions/workflows/job-collector.yml/badge.svg?branch=master)](https://github.com/margostino/climateline-processor/actions/workflows/job-collector.yml)
[![publisher-dispatcher](https://github.com/margostino/climateline-processor/actions/workflows/publisher-dispatcher.yml/badge.svg?branch=master)](https://github.com/margostino/climateline-processor/actions/workflows/publisher-dispatcher.yml)

This repo implements the job to fetch the daily news of Climate Change.  
The job retrieve every news into the Admin Telegram Bot in order to approve, edit and push new article into [Climateline](https://climateline.vercel.app/) website.

The main motivation behind is to automate and easy the news uploading with zero line of code from everywhere and whenever you want using a Telegram Bot connection. 

### Architecture

![](assets/architecture.png#100x)

### Features

- [x] Cron job for fetching news
- [x] Telegram bot for handling news workflow (upload, update) towards [Climateline](https://climateline.vercel.app/)
- [x] Basic and short-time caching
- [x] Multi-Source News Fetcher
- [ ] Create entirely new articles
- [ ] Detect duplicated article before pushing
- [ ] Pushing by replying
- [ ] Detect automatically properties such as location, source, category
- [ ] Deploy only if build and tests pass
- [ ] Isolate news fetcher API

