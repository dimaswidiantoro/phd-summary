---
layout: default
title: Jekyll Installation Guide
---

Certainly! Here is a comprehensive documentation that compiles all the issues you faced during the Jekyll installation and layout setup, along with the solutions to resolve them.

---

# Troubleshooting Jekyll Installation and Layout Issues on macOS

## Introduction

This guide documents common issues encountered when installing Jekyll on macOS, particularly the `Gem::FilePermissionError` and compatibility issues with the `google-protobuf` gem. It also covers troubleshooting steps for resolving 404 errors and configuring the layout properly. Follow these step-by-step solutions to resolve these issues.

## Table of Contents

1. [Problem Overview](#problem-overview)
2. [Installing Jekyll and Bundler](#installing-jekyll-and-bundler)
   - [Using `sudo` to Install Gems](#using-sudo-to-install-gems)
   - [Using a Ruby Version Manager (`rbenv`)](#using-a-ruby-version-manager-rbenv)
3. [Creating a New Jekyll Site](#creating-a-new-jekyll-site)
4. [Editing Jekyll Configuration and Content](#editing-jekyll-configuration-and-content)
5. [Serving the Jekyll Site Locally](#serving-the-jekyll-site-locally)
6. [Pushing Changes to GitHub](#pushing-changes-to-github)
7. [Troubleshooting 404 Errors](#troubleshooting-404-errors)
8. [Customizing Layouts and Adding Navigation](#customizing-layouts-and-adding-navigation)
9. [Conclusion](#conclusion)

## Problem Overview

### Issues Faced

1. **Gem::FilePermissionError**:
   - Error: `You don't have write permissions for the /Library/Ruby/Gems/2.6.0 directory.`
   - Cause: Insufficient permissions to install gems system-wide.

2. **Compatibility Issues with `google-protobuf` Gem**:
   - Error: `There are no versions of google-protobuf (>= 3.25, < 5.0) compatible with your Ruby & RubyGems.`
   - Cause: Outdated Ruby version incompatible with required gems.

3. **404 Errors**:
   - Error: `404 Page not found :(`
   - Cause: Incorrect file paths or misconfigured `_config.yml`.

## Installing Jekyll and Bundler

### Using `sudo` to Install Gems

1. **Install Jekyll and Bundler with Elevated Permissions**:
   ```sh
   sudo gem install jekyll bundler
   ```

### Using a Ruby Version Manager (`rbenv`)

1. **Install `rbenv` and `ruby-build`**:
   ```sh
   brew install rbenv ruby-build
   ```

2. **Set Up `rbenv`**:
   - Add `rbenv` to your shell configuration file:
     ```sh
     echo 'eval "$(rbenv init -)"' >> ~/.zshrc
     source ~/.zshrc
     ```
   - For `bash`, add it to `~/.bashrc` or `~/.bash_profile`.

3. **Install a Newer Version of Ruby**:
   ```sh
   rbenv install 3.1.0
   rbenv global 3.1.0
   ```

4. **Rehash `rbenv`**:
   ```sh
   rbenv rehash
   ```

5. **Verify Ruby Installation**:
   ```sh
   ruby -v
   ```

6. **Install Jekyll and Bundler**:
   ```sh
   gem install jekyll bundler
   ```

## Creating a New Jekyll Site

1. **Create a New Jekyll Site**:
   ```sh
   jekyll new my-documentation
   cd my-documentation
   ```

2. **Verify the Presence of the `Gemfile`**:
   ```sh
   ls
   ```

3. **Install Dependencies**:
   ```sh
   bundle install
   ```

## Editing Jekyll Configuration and Content

1. **Edit `_config.yml`**:
   Ensure your `_config.yml` has the correct settings:

   ```yaml
   title: My Documentation Site
   email: your-email@example.com
   description: >- # this means to ignore newlines until "baseurl:"
     A documentation site for my project.
   baseurl: "" # the subpath of your site, e.g. /blog
   url: "http://localhost:4000" # the base hostname & protocol for your site

   theme: minima

   markdown: kramdown
   highlighter: rouge

   kramdown:
     input: GFM
     syntax_highlighter_opts:
       line_numbers: true
   ```

2. **Create and Edit Markdown Files**:
   - Create a new Markdown file:
     ```sh
     nano jekyll-installation-guide.md
     ```
   - Add content to `jekyll-installation-guide.md`:
     ```markdown
     ---
     layout: default
     title: Jekyll Installation Guide
     ---

     # Troubleshooting Jekyll Installation on macOS

     This guide documents common issues encountered when installing Jekyll on macOS, particularly the `Gem::FilePermissionError` and compatibility issues with the `google-protobuf` gem. It provides step-by-step solutions to resolve these issues using a Ruby version manager (`rbenv`).

     ## Table of Contents

     1. [Problem Overview](#problem-overview)
     2. [Installing Jekyll and Bundler](#installing-jekyll-and-bundler)
        - [Using `sudo` to Install Gems](#using-sudo-to-install-gems)
        - [Using a Ruby Version Manager (`rbenv`)](#using-a-ruby-version-manager-rbenv)
     3. [Creating a New Jekyll Site](#creating-a-new-jekyll-site)
     4. [Editing Jekyll Configuration and Content](#editing-jekyll-configuration-and-content)
     5. [Serving the Jekyll Site Locally](#serving-the-jekyll-site-locally)
     6. [Pushing Changes to GitHub](#pushing-changes-to-github)
     7. [Troubleshooting 404 Errors](#troubleshooting-404-errors)
     8. [Customizing Layouts and Adding Navigation](#customizing-layouts-and-adding-navigation)
     9. [Conclusion](#conclusion)
     ```

## Serving the Jekyll Site Locally

1. **Serve the Site Locally**:
   ```sh
   bundle exec jekyll serve
   ```

2. **Preview the Site**:
   Open your browser and navigate to `http://localhost:4000`.

## Pushing Changes to GitHub

1. **Commit and Push Changes to GitHub**:
   ```sh
   git add .
   git commit -m "Added Jekyll installation guide"
   git push origin main
   ```

## Troubleshooting 404 Errors

1. **Verify the File Structure**:
   Ensure that the files are in the correct directories and named correctly.

   ```sh
   ls
   ```

2. **Check `_config.yml` Configuration**:
   Ensure your `_config.yml` has the correct base URL settings:

   ```yaml
   baseurl: "" # the subpath of your site, e.g., /blog
   url: "http://localhost:4000" # the base hostname & protocol for your site
   ```

3. **Verify Front Matter in Markdown Files**:
   Ensure that each Markdown file has the correct front matter at the top. For example, `jekyll-installation-guide.md` should start with:

   ```markdown
   ---
   layout: default
   title: Jekyll Installation Guide
   ---
   ```

4. **Check Links in `index.md`**:
   Ensure that the links in your `index.md` file are correct. For example:

   ```markdown
   ---
   layout: default
   title: Home
   ---

   # Welcome to My Documentation Site

   This is the homepage of my documentation site. Use the navigation links below to find what you need.

   - [Jekyll Installation Guide](./jekyll-installation-guide.html)
   - [MongoDB Documentation](./mongodb.html)
   - [Go Documentation](./go.html)
   - [React Documentation](./react.html)
   ```

## Customizing Layouts and Adding Navigation

1. **Create a Custom Layout**:

   ```sh
   mkdir -p _layouts
   nano _layouts/documentation.html
   ```

   Add the following content to `documentation.html`:

   ```html
   <!DOCTYPE html>
   <html>
   <head>
       <meta charset="utf-8">
       <title>{{ page.title }}</title>
       <link rel="stylesheet" href="{{ "/assets/css/style.css" | relative_url }}">
   </head>
   <body>
       <header>
           <h1>{{ site.title }}</h1>
       </header>
       <nav>
           <ul>
               <li><a href="{{ "/" | relative_url }}">Home</a></li>
               <li><a href="{{ "/jekyll-installation-guide" | relative_url }}">Jekyll Installation Guide</a></li>
               <li><a href="{{ "/mongodb" | relative_url }}">MongoDB</a></li>
               <li><a href="{{ "/go" | relative_url }}">Go</a></li>
               <li><a href="{{ "/react"

 | relative_url }}">React</a></li>
           </ul>
       </nav>
       <main>
           {{ content }}
       </main>
       <footer>
           <p>&copy; {{ site.time | date: "%Y" }} {{ site.title }}. All rights reserved.</p>
       </footer>
   </body>
   </html>
   ```

2. **Link Custom CSS in Layout**:

   Create a custom CSS file (e.g., `assets/css/style.css`) and link it in your layout.

   ```sh
   mkdir -p assets/css
   nano assets/css/style.css
   ```

   Add custom styles:

   ```css
   /* Custom styles for code blocks */
   pre {
     background-color: #2d2d2d;
     color: #f8f8f2;
     padding: 10px;
     border-radius: 5px;
     overflow: auto;
   }

   code {
     background-color: #2d2d2d;
     color: #f8f8f2;
     padding: 2px 5px;
     border-radius: 3px;
   }
   ```

   Edit your layout file to include the custom CSS:

   ```sh
   nano _layouts/default.html
   ```

   Add the following line in the `<head>` section:

   ```html
   <link rel="stylesheet" href="{{ "/assets/css/style.css" | relative_url }}">
   ```

## Conclusion

By following this guide, you can resolve common issues encountered when installing Jekyll on macOS, create a new Jekyll site, and add your documentation content. This guide also demonstrates how to serve your Jekyll site locally and push changes to GitHub. Additionally, it provides steps to troubleshoot 404 errors and customize your site layout.

---

You can use this documentation on your Jekyll site by creating a new Markdown file (e.g., `jekyll-installation-guide.md`) and adding it to your site. This guide should help others facing similar issues with Jekyll installation and configuration. If you have any further questions or need additional assistance, feel free to ask!