# Caffeine

[RU](README.ru.md)

[![Go Coverage](https://github.com/akovardin/adscoffee/wiki/coverage.svg)](https://raw.githack.com/wiki/akovardin/adscoffee/coverage.html)
Here is the translation preserving the original formatting and markdown structure.

## Open-Source Ad Platform for Developers

Caffeine is a full-featured open-source advertising platform built to give developers, publishers, and open-source communities complete control over monetizing their projects. We offer a transparent, customizable, and ethical alternative to major ad networks.

The philosophy of Caffeine is built on four fundamental principles.

Transparency is the foundation of trust. The platform's entire source code is open for independent audit. You will always know exactly how the ad selection algorithm works, what data is collected, and how it is processed. No hidden algorithms or "black boxes."

Control is entirely in your hands. You independently manage all aspects of monetization: from the design and placement of ad units to the fine-tuning of targeting and ad frequency. The platform is your tool, not a set of imposed rules.

Ethics is at the core of our approach to the user. The system is designed with respect for the audience. You can easily set strict rules for ad quality and relevance, avoid intrusive formats, and guaranteed comply with regulatory standards such as GDPR.

Community is our key audience and goal. Caffeine is specifically created to support the open-knowledge ecosystem. The platform is ideally suited for monetizing documentation, GitHub projects, personal blogs, and niche applications, areas where large commercial networks are often ineffective or excessive.

## Who is Caffeine For?

1.  Open-Source Developers: Monetize documentation pages, demo sites, or repositories without breaking community trust.
2.  Small Publishers and Bloggers: Get an alternative to AdSense with greater flexibility and no entry barriers.
3.  Game and App Developers (Web, Mobile, Desktop): Use it as a primary ad network or as your own system for direct advertisers.
4.  Corporate and Internal Projects: Display relevant internal announcements (about new services, events) to employees or partners.
5.  Researchers and Students: A perfect testbed for studying targeting algorithms, mediation, and ad campaign analysis.

## Current Functionality (Core Features)

Core System:

- Modular Plugin-Based Architecture: Easily extend functionality (add new formats, data sources, analytics systems) without changing the core.
- Ad Server: Accepts requests from clients and returns the most relevant ad.
- Administrative Panel: A user-friendly interface for managing ad campaigns and creatives.
- Analytics and Reporting based on ClickHouse and Redash.

Integration Clients (SDK / Libraries):

- Web Client: A lightweight JS library for embedding ads on websites.
- Client for Native Applications: A universal client adaptable for mobile (iOS/Android) and desktop (Windows, macOS, Linux) applications.

## Technology Stack

- Backend: Go (high performance, static typing for reliability).
- Frontend (Admin Panel): Uses qor5 with custom improvements.
- Data Storage: PostgreSQL, ClickHouse, Redis.
- Client Libraries: JavaScript (Web) with the possibility of creating wrappers for other languages.

## Why You Should Consider Caffeine

- Eliminate the "Black Box": You see all the logic of operation, instead of trusting unclear algorithms.
- Flexibility: Adapt the platform to your unique needs, be it a special ad format or integration with an internal CRM.
- Savings: Don't pay huge commissions to intermediaries (often 30-50%). Pay only for hosting and infrastructure.
- User Trust: Honest advertising and transparency increase audience loyalty.
- Future Development: You participate in the platform's evolution. Your pull requests and ideas shape the project's future.

## What is Currently Implemented

- [x] Plugin System
- [ ] Admin Panel for Ad Configuration
- [ ] Ad Selection Server
- [ ] Analytics and Statistics Display
- [ ] Client for Displaying Ads in Applications
- [ ] Client for Displaying Ads on Websites
- [ ] Ad Mediation

Caffeine is more than just a tool. It is an attempt to make the advertising ecosystem healthier, more open, and more beneficial for those who create real value on the internet.