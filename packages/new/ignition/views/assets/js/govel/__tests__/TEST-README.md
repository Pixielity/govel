# GoVel Modules Test Suite

This directory contains a comprehensive test suite to verify that all GoVel modules work correctly with RequireJS AMD dependencies.

## Architecture

All modules now use AMD (Asynchronous Module Definition) format via RequireJS:

### AMD Modules Structure:
```
govel/
├── constants/
│   └── config.constant.js    (AMD module: returns CONFIG)
├── utils/
│   ├── dom.util.js          (AMD module: returns DOM)  
│   └── logger.util.js       (AMD module: returns Logger)
├── govel-ui.js              (AMD module: depends on config, logger, dom)
├── govel.js                 (AMD module: depends on govel-ui)
└── test.html                (Test page with RequireJS)
```

### Dependencies:
- `govel.js` depends on `govel-ui.js`
- `govel-ui.js` depends on `constants/config.constant`, `utils/logger.util`, `utils/dom.util`
- All constants and utils are now AMD modules

## Running Tests

1. **Open test.html in a web browser**
   ```bash
   open test.html
   # or
   python -m http.server 8000  # then visit http://localhost:8000/test.html
   ```

2. **Use the test buttons:**
   - 🧪 **Run All Tests** - Runs complete test suite
   - 🌍 **Test Global Modules** - Info about module architecture  
   - 📦 **Test AMD Modules** - Tests individual AMD modules
   - 🔗 **Test Integration** - Tests module interactions
   - 🧹 **Clear Results** - Clears test output

## Test Coverage

### AMD Module Tests:
- ✅ CONFIG module loading and structure
- ✅ Logger module functionality  
- ✅ DOM utility functions
- ✅ GovelIgnitionUI class instantiation

### Integration Tests:
- ✅ Main module (govel.js) loading
- ✅ Dependency chain resolution
- ✅ Module interactions and data flow
- ✅ Mock DOM operations with CONFIG selectors

### Expected Results:
- All modules should load successfully via RequireJS
- No global variables needed (everything is AMD)
- Dependencies resolve automatically
- Mock DOM tests verify selector functionality

## Troubleshooting

### Common Issues:

1. **Module not found**: Check file paths in `require()` calls
2. **Dependency errors**: Ensure all AMD dependencies are correctly specified  
3. **RequireJS 404s**: Make sure RequireJS CDN is accessible

### Debug Mode:
Add `?debug=1` to the URL to enable debug logging in the GoVel modules.

## Module Usage Example

```javascript
// Load modules via RequireJS
require(['govel-ui', 'constants/config.constant'], function(GovelIgnitionUI, CONFIG) {
    var ui = new GovelIgnitionUI({ debug: true });
    console.log('Config selectors:', CONFIG.SELECTORS);
    ui.init();
});
```

## Files Overview

- **test.html** - Main test page with RequireJS setup
- **test-suite.js** - Comprehensive test framework
- **constants/config.constant.js** - Configuration constants (AMD)
- **utils/logger.util.js** - Logging utility (AMD)
- **utils/dom.util.js** - DOM manipulation utilities (AMD) 
- **govel-ui.js** - Main UI customization class (AMD)
- **govel.js** - Entry point module (AMD)
