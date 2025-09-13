/**
 * GoVel Modules Test Suite
 * 
 * Comprehensive test suite to verify that all GoVel modules work correctly
 * with RequireJS AMD dependencies.
 * 
 * @version 1.0.0
 * @author GoVel Team
 * @created 2025-09-11
 */

// Test framework
var TestFramework = {
    results: [],
    currentSection: null,
    
    /**
     * Add a test result
     */
    addResult: function(name, passed, message, details) {
        this.results.push({
            section: this.currentSection,
            name: name,
            passed: passed,
            message: message,
            details: details || null,
            timestamp: new Date().toISOString()
        });
        this.displayResult(name, passed, message, details);
    },
    
    /**
     * Start a new test section
     */
    startSection: function(sectionName) {
        this.currentSection = sectionName;
        var container = document.getElementById('test-results');
        if (!container) {
            console.error('Test results container not found');
            return;
        }
        var sectionDiv = document.createElement('div');
        sectionDiv.className = 'test-container';
        var sectionId = sectionName ? sectionName.replace(/\s+/g, '-').toLowerCase() : 'default';
        sectionDiv.innerHTML = '<div class="section-header"><h2>' + sectionName + '</h2></div><div id="section-' + sectionId + '"></div>';
        container.appendChild(sectionDiv);
    },
    
    /**
     * Display a test result
     */
    displayResult: function(name, passed, message, details) {
        var sectionId = 'section-' + (this.currentSection ? this.currentSection.replace(/\s+/g, '-').toLowerCase() : 'default');
        var section = document.getElementById(sectionId);
        
        // If section doesn't exist, create a default one
        if (!section) {
            this.startSection(this.currentSection || 'Default Tests');
            section = document.getElementById(sectionId);
        }
        
        if (!section) {
            console.error('Could not create or find test section');
            return;
        }
        
        var resultDiv = document.createElement('div');
        resultDiv.className = 'test-result ' + (passed ? 'test-pass' : 'test-fail');
        
        var icon = passed ? '✓' : '✗';
        var content = icon + ' ' + name + ': ' + message;
        
        if (details) {
            content += '<br><pre>' + JSON.stringify(details, null, 2) + '</pre>';
        }
        
        resultDiv.innerHTML = content;
        section.appendChild(resultDiv);
    },
    
    /**
     * Show test summary
     */
    showSummary: function() {
        var passed = this.results.filter(r => r.passed).length;
        var failed = this.results.filter(r => !r.passed).length;
        var total = this.results.length;
        
        var container = document.getElementById('test-results');
        if (!container) {
            console.error('Test results container not found for summary');
            return { passed: passed, failed: failed, total: total };
        }
        
        var summaryDiv = document.createElement('div');
        summaryDiv.className = 'test-summary ' + (failed === 0 ? 'test-pass' : 'test-fail');
        summaryDiv.innerHTML = 'Tests: ' + passed + '/' + total + ' passed' + (failed > 0 ? ', ' + failed + ' failed' : '');
        container.insertBefore(summaryDiv, container.firstChild);
        
        return { passed: passed, failed: failed, total: total };
    },
    
    /**
     * Clear all results
     */
    clear: function() {
        this.results = [];
        this.currentSection = null;
        var container = document.getElementById('test-results');
        if (container) {
            container.innerHTML = '';
        }
    }
};

/**
 * Ensure DOM is ready before running tests
 */
function ensureReady(callback) {
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', callback);
    } else {
        callback();
    }
}

/**
 * Test global functions
 */
function runAllTests() {
    ensureReady(function() {
        TestFramework.clear();
        runAMDModuleTests();
        setTimeout(function() {
            runIntegrationTests();
            TestFramework.showSummary();
        }, 500);
    });
}

function runGlobalTests() {
    ensureReady(function() {
        TestFramework.clear();
        TestFramework.startSection('Global Tests');
        TestFramework.addResult('Global Test Info', true, 'All modules are now AMD-based, no global tests needed');
        TestFramework.showSummary();
    });
}

function runAMDTests() {
    ensureReady(function() {
        TestFramework.clear();
        runAMDModuleTests();
        setTimeout(function() {
            TestFramework.showSummary();
        }, 1000);
    });
}

function runIntegrationTests() {
    ensureReady(function() {
        if (TestFramework.results.length === 0) {
            TestFramework.clear();
        }
        runIntegrationTestsOnly();
        setTimeout(function() {
            TestFramework.showSummary();
        }, 1500);
    });
}

function testMockDOM() {
    ensureReady(function() {
        TestFramework.clear();
        TestFramework.startSection('Mock DOM Tests');
        
        require(['../govel-ui', '../constants/config.constant', '../utils/dom.util'], 
        function(GovelIgnitionUI, CONFIG, DOM) {
            try {
                // Test that mock elements exist
                var mockContainer = document.querySelector('.mock-ignition');
                TestFramework.addResult('Mock Container', 
                    mockContainer !== null,
                    mockContainer ? 'Mock container found' : 'Mock container not found'
                );
                
                if (mockContainer) {
                    // Test individual selectors
                    var footer = DOM.select(CONFIG.SELECTORS.FOOTER, mockContainer);
                    TestFramework.addResult('Footer Selector', 
                        footer !== null,
                        footer ? 'Footer element found: ' + footer.textContent : 'Footer element not found'
                    );
                    
                    var phpDocsLink = DOM.select(CONFIG.SELECTORS.PHP_DOCS_LINK, mockContainer);
                    TestFramework.addResult('PHP Docs Link', 
                        phpDocsLink !== null,
                        phpDocsLink ? 'PHP docs link found: ' + phpDocsLink.href : 'PHP docs link not found'
                    );
                    
                    var flareLink = DOM.select(CONFIG.SELECTORS.FLARE_LINK, mockContainer);
                    TestFramework.addResult('Flare Link', 
                        flareLink !== null,
                        flareLink ? 'Flare link found: ' + flareLink.href : 'Flare link not found'
                    );
                    
                    var ignitionLogo = DOM.select(CONFIG.SELECTORS.IGNITION_LOGO, mockContainer);
                    TestFramework.addResult('Ignition Logo', 
                        ignitionLogo !== null,
                        ignitionLogo ? 'Ignition logo found' : 'Ignition logo not found'
                    );
                    
                    // Test GoVel UI transformations
                    TestFramework.addResult('Starting Transformations', true, 'Testing GoVel UI transformations on mock elements');
                    
                    var ui = new GovelIgnitionUI({ debug: true });
                    
                    // Test individual transformation methods
                    ui.hideFooter().then(function(result) {
                        TestFramework.addResult('Hide Footer Transform', 
                            result === true,
                            result ? 'Footer hidden successfully' : 'Footer hide failed'
                        );
                    }).catch(function(error) {
                        TestFramework.addResult('Hide Footer Transform', false, 'Footer hide error: ' + error.message);
                    });
                    
                    ui.hideFlareMenuItem().then(function(result) {
                        TestFramework.addResult('Hide Flare MenuItem Transform', 
                            result === true,
                            result ? 'Flare menu item hidden successfully' : 'Flare menu item hide failed'
                        );
                    }).catch(function(error) {
                        TestFramework.addResult('Hide Flare MenuItem Transform', false, 'Flare menu item hide error: ' + error.message);
                    });
                    
                    ui.updatePhpToGo().then(function(result) {
                        TestFramework.addResult('PHP to Go Transform', 
                            result === true,
                            result ? 'PHP to Go transformation successful' : 'PHP to Go transformation failed'
                        );
                    }).catch(function(error) {
                        TestFramework.addResult('PHP to Go Transform', false, 'PHP to Go transformation error: ' + error.message);
                    });
                }
                
            } catch (error) {
                TestFramework.addResult('Mock DOM Test Setup', false, 'Test setup failed: ' + error.message);
            }
        }, function(error) {
            TestFramework.addResult('Mock DOM Module Load', false, 'Failed to load modules: ' + error.message);
        });
        
        setTimeout(function() {
            TestFramework.showSummary();
        }, 2000);
    });
}

function clearResults() {
    ensureReady(function() {
        TestFramework.clear();
    });
}

/**
 * Test AMD modules
 */
function runAMDModuleTests() {
    TestFramework.startSection('AMD Module Tests');
    
    // Test CONFIG module
    require(['../constants/config.constant'], function(CONFIG) {
        try {
            TestFramework.addResult('CONFIG Module Load', true, 'CONFIG module loaded successfully');
            
            // Test CONFIG properties
            TestFramework.addResult('CONFIG.TIMING', 
                CONFIG.TIMING && typeof CONFIG.TIMING === 'object', 
                CONFIG.TIMING ? 'TIMING configuration exists' : 'TIMING configuration missing',
                CONFIG.TIMING
            );
            
            TestFramework.addResult('CONFIG.URL_MAPPINGS', 
                CONFIG.URL_MAPPINGS && typeof CONFIG.URL_MAPPINGS === 'object', 
                CONFIG.URL_MAPPINGS ? 'URL_MAPPINGS configuration exists' : 'URL_MAPPINGS configuration missing',
                Object.keys(CONFIG.URL_MAPPINGS || {})
            );
            
            TestFramework.addResult('CONFIG.SELECTORS', 
                CONFIG.SELECTORS && typeof CONFIG.SELECTORS === 'object', 
                CONFIG.SELECTORS ? 'SELECTORS configuration exists' : 'SELECTORS configuration missing',
                Object.keys(CONFIG.SELECTORS || {})
            );
            
        } catch (error) {
            TestFramework.addResult('CONFIG Module', false, 'CONFIG module test failed: ' + error.message);
        }
    }, function(error) {
        TestFramework.addResult('CONFIG Module Load', false, 'Failed to load CONFIG module: ' + error.message);
    });
    
    // Test Logger module
    require(['../utils/logger.util'], function(Logger) {
        try {
            TestFramework.addResult('Logger Module Load', true, 'Logger module loaded successfully');
            
            var logger = new Logger(true);
            TestFramework.addResult('Logger Instance', 
                logger && typeof logger === 'object', 
                'Logger instance created successfully'
            );
            
            TestFramework.addResult('Logger Methods', 
                typeof logger.info === 'function' && typeof logger.warn === 'function' && typeof logger.error === 'function',
                'Logger has required methods (info, warn, error)'
            );
            
        } catch (error) {
            TestFramework.addResult('Logger Module', false, 'Logger module test failed: ' + error.message);
        }
    }, function(error) {
        TestFramework.addResult('Logger Module Load', false, 'Failed to load Logger module: ' + error.message);
    });
    
    // Test DOM module
    require(['../utils/dom.util'], function(DOM) {
        try {
            TestFramework.addResult('DOM Module Load', true, 'DOM module loaded successfully');
            
            TestFramework.addResult('DOM Methods', 
                typeof DOM.select === 'function' && 
                typeof DOM.selectAll === 'function' && 
                typeof DOM.isVisible === 'function' &&
                typeof DOM.retry === 'function',
                'DOM has required methods (select, selectAll, isVisible, retry)'
            );
            
            // Test DOM.select
            var body = DOM.select('body');
            TestFramework.addResult('DOM.select', 
                body && body.tagName === 'BODY', 
                'DOM.select can find body element'
            );
            
        } catch (error) {
            TestFramework.addResult('DOM Module', false, 'DOM module test failed: ' + error.message);
        }
    }, function(error) {
        TestFramework.addResult('DOM Module Load', false, 'Failed to load DOM module: ' + error.message);
    });
    
    // Test GovelIgnitionUI module
    require(['../govel-ui'], function(GovelIgnitionUI) {
        try {
            TestFramework.addResult('GovelIgnitionUI Module Load', true, 'GovelIgnitionUI module loaded successfully');
            
            TestFramework.addResult('GovelIgnitionUI Constructor', 
                typeof GovelIgnitionUI === 'function', 
                'GovelIgnitionUI is a constructor function'
            );
            
            var ui = new GovelIgnitionUI({ debug: true });
            TestFramework.addResult('GovelIgnitionUI Instance', 
                ui && typeof ui === 'object', 
                'GovelIgnitionUI instance created successfully'
            );
            
            TestFramework.addResult('GovelIgnitionUI Methods', 
                typeof ui.init === 'function' && 
                typeof ui.getStatus === 'function',
                'GovelIgnitionUI has required methods (init, getStatus)'
            );
            
        } catch (error) {
            TestFramework.addResult('GovelIgnitionUI Module', false, 'GovelIgnitionUI module test failed: ' + error.message);
        }
    }, function(error) {
        TestFramework.addResult('GovelIgnitionUI Module Load', false, 'Failed to load GovelIgnitionUI module: ' + error.message);
    });
}

/**
 * Test integration between modules
 */
function runIntegrationTestsOnly() {
    TestFramework.startSection('Integration Tests');
    
    // Test full module integration
    require(['../govel'], function(GovelMain) {
        try {
            TestFramework.addResult('Main Module Load', true, 'Main govel module loaded successfully');
            
            TestFramework.addResult('Main Module Structure', 
                GovelMain && typeof GovelMain.init === 'function',
                'Main module has init function'
            );
            
        } catch (error) {
            TestFramework.addResult('Main Module', false, 'Main module test failed: ' + error.message);
        }
    }, function(error) {
        TestFramework.addResult('Main Module Load', false, 'Failed to load main module: ' + error.message);
    });
    
    // Test dependency chain
    require(['../govel-ui', '../constants/config.constant', '../utils/logger.util', '../utils/dom.util'], 
    function(GovelIgnitionUI, CONFIG, Logger, DOM) {
        try {
            TestFramework.addResult('Dependency Chain', true, 'All dependencies loaded together successfully');
            
            // Test that GovelIgnitionUI can use the dependencies
            var ui = new GovelIgnitionUI({ debug: true });
            var status = ui.getStatus();
            
            TestFramework.addResult('Module Interaction', 
                status && typeof status === 'object' && typeof status.initialized === 'boolean',
                'GovelIgnitionUI can interact with dependencies'
            );
            
            // Test CONFIG usage
            TestFramework.addResult('CONFIG Integration', 
                CONFIG.SELECTORS && CONFIG.SELECTORS.FOOTER,
                'CONFIG constants are accessible',
                { footer_selector: CONFIG.SELECTORS.FOOTER }
            );
            
            // Test Logger usage
            var logger = new Logger(true);
            TestFramework.addResult('Logger Integration', 
                logger && typeof logger.info === 'function',
                'Logger utility is functional'
            );
            
            // Test DOM usage
            var testElement = DOM.select('body');
            TestFramework.addResult('DOM Integration', 
                testElement && testElement.tagName === 'BODY',
                'DOM utility is functional'
            );
            
        } catch (error) {
            TestFramework.addResult('Dependency Chain', false, 'Dependency integration test failed: ' + error.message);
        }
    }, function(error) {
        TestFramework.addResult('Dependency Chain Load', false, 'Failed to load dependency chain: ' + error.message);
    });
    
    // Test mock DOM operations  
    setTimeout(function() {
        require(['../govel-ui', '../constants/config.constant', '../utils/dom.util'], 
        function(GovelIgnitionUI, CONFIG, DOM) {
            try {
                var mockContainer = document.querySelector('.mock-ignition');
                if (mockContainer) {
                    // Test selector functionality
                    var footer = DOM.select(CONFIG.SELECTORS.FOOTER, mockContainer);
                    TestFramework.addResult('Selector Test', 
                        footer !== null,
                        footer ? 'Footer selector works with CONFIG' : 'Footer selector failed'
                    );
                    
                    var phpLink = DOM.select(CONFIG.SELECTORS.PHP_DOCS_LINK, mockContainer);
                    TestFramework.addResult('PHP Link Selector', 
                        phpLink !== null,
                        phpLink ? 'PHP docs link selector works' : 'PHP docs link selector failed'
                    );
                }
                
            } catch (error) {
                TestFramework.addResult('Mock DOM Test', false, 'Mock DOM test failed: ' + error.message);
            }
        });
    }, 1000);
}
