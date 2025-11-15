# Viewing Documentation

This project uses [Docsify](https://docsify.js.org/) for documentation with integrated Swagger UI for API documentation.

## 📚 Online Documentation

The documentation is hosted on GitHub Pages:

**https://ericomonteiro.github.io/easy-queue-go/**

## 💻 Local Documentation

### Option 1: Using Python HTTP Server

Navigate to the docs directory and start a local server:

```bash
cd docs
python3 -m http.server 3000
```

Then open your browser at: **http://localhost:3000**

### Option 2: Using Node.js HTTP Server

Install `http-server` globally (one time only):

```bash
npm install -g http-server
```

Navigate to the docs directory and start the server:

```bash
cd docs
http-server -p 3000
```

Then open your browser at: **http://localhost:3000**

### Option 3: Using Docsify CLI

Install Docsify CLI (one time only):

```bash
npm install -g docsify-cli
```

Navigate to the docs directory and start the server:

```bash
cd docs
docsify serve
```

Then open your browser at: **http://localhost:3000**

### Option 4: Using VS Code Live Server Extension

1. Install the "Live Server" extension in VS Code
2. Right-click on `docs/index.html`
3. Select "Open with Live Server"

## 📖 Documentation Structure

- **Home** - Project overview and quick start
- **Getting Started** - Setup instructions
- **Project Structure** - Code organization
- **Database** - Schema and migrations
- **Features** - Feature documentation
- **API** - API documentation with Swagger UI
- **Product** - Product vision and roadmap

## 🔄 Updating Swagger Documentation

After making changes to API handlers:

1. Generate new Swagger files:
```bash
make swagger-generate
```

2. The following files will be updated:
   - `docs/docs.go` - Go code (gitignored)
   - `docs/swagger.json` - OpenAPI spec in JSON
   - `docs/swagger.yaml` - OpenAPI spec in YAML

3. Refresh your browser to see the updated API documentation

## 🌐 Swagger UI Access

You can access the Swagger UI in two ways:

1. **Integrated in Docsify** - Navigate to API > Swagger UI in the documentation
2. **Standalone** - When the Go application is running: http://localhost:8080/swagger/index.html

## 📝 Notes

- The documentation uses Docsify's sidebar navigation
- Swagger UI is embedded directly in the documentation
- All API endpoints can be tested directly from the browser
- Authentication is supported via the "Authorize" button in Swagger UI
