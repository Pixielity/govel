/**
 * GoVel Ignition UI Customization Module
 * 
 * This module transforms the Laravel Ignition error page interface to be Go-specific.
 * It provides seamless PHP-to-Go branding conversion, UI element customization,
 * and dynamic content modification for Go-based error reporting.
 * 
 * @version 1.0.0
 * @author GoVel Team
 * @created 2025-09-11
 * @requires RequireJS
 * @module govel-ui
 */

define(['./constants/config.constant', './utils/logger.util', './utils/dom.util'], function(CONFIG, Logger, DOM) {
    'use strict';
    
    /**
     * Main GoVel Ignition UI customization class
     * @class GovelIgnitionUI
     */
    function GovelIgnitionUI(options) {
        this.options = options || {};
        this.logger = new Logger(this.options.debug);
        this.initialized = false;
        this.modifications = new Set();
        
        // Store reference for instance methods
        var self = this;
        
        // Store logger instance globally for static methods
        Logger.instance = this.logger;
    }

    /**
     * Initialize all UI modifications
     * @returns {Promise<void>}
     */
    GovelIgnitionUI.prototype.init = async function() {
        if (this.initialized) {
            this.logger.warn('UI customizer already initialized');
            return;
        }

        this.logger.info('Initializing GoVel Ignition UI customizer');

        try {
            // Wait for DOM to be ready
            await new Promise(resolve => DOM.ready(resolve));
            
            // Add small delay to ensure all content is loaded
            await new Promise(resolve => 
                setTimeout(resolve, this.options.initDelay || CONFIG.TIMING.DOM_READY_DELAY)
            );

            // Run all modifications
            await this.runAllModifications();
            
            this.initialized = true;
            this.logger.info('GoVel Ignition UI customizer initialized successfully');
            
        } catch (error) {
            this.logger.error('Failed to initialize UI customizer:', error);
            throw error;
        }
    };

    /**
     * Execute all UI modifications
     * @private
     * @returns {Promise<void>}
     */
    GovelIgnitionUI.prototype.runAllModifications = async function() {
        const modifications = [
            { name: 'hideFooter', fn: this.hideFooter.bind(this) },
            { name: 'hideFlareMenuItem', fn: this.hideFlareMenuItem.bind(this) },
            { name: 'updatePhpToGo', fn: this.updatePhpToGo.bind(this) },
            { name: 'updateLanguageVersion', fn: this.updateLanguageVersion.bind(this) },
            { name: 'updateSourceLinks', fn: this.updateSourceLinks.bind(this) },
            { name: 'addGoGopherIcon', fn: this.addGoGopherIcon.bind(this) }
        ];

        for (var i = 0; i < modifications.length; i++) {
            var modification = modifications[i];
            try {
                await modification.fn();
                this.modifications.add(modification.name);
                this.logger.info('✓ ' + modification.name + ' completed');
            } catch (error) {
                this.logger.error('✗ ' + modification.name + ' failed:', error);
            }
        }

        this.logger.info('Modifications completed: ' + this.modifications.size + '/' + modifications.length);
    };

    /**
     * Hide the Laravel/Flare footer branding
     * Removes footer elements that contain Laravel/Flare branding to provide
     * a clean Go-specific error page experience.
     * 
     * @returns {Promise<boolean>} Success status
     */
    GovelIgnitionUI.prototype.hideFooter = async function() {
        return DOM.retry(() => {
            const footer = DOM.select(CONFIG.SELECTORS.FOOTER);
            if (footer && DOM.isVisible(footer)) {
                footer.style.display = 'none';
                this.logger.info('Footer hidden successfully');
                return true;
            }
            return false;
        });
    };

    /**
     * Hide the Flare menu item (third li element)
     * Removes the Flare link from the top navigation menu to maintain
     * Go-specific branding without Flare references.
     * 
     * @returns {Promise<boolean>} Success status
     */
    GovelIgnitionUI.prototype.hideFlareMenuItem = async function() {
        return DOM.retry(() => {
            const flareLink = DOM.select(CONFIG.SELECTORS.FLARE_LINK);
            if (flareLink && DOM.isVisible(flareLink)) {
                // Find the parent <li> element
                const flareMenuItem = flareLink.closest('li');
                if (flareMenuItem) {
                    flareMenuItem.style.display = 'none';
                    this.logger.info('Flare menu item hidden successfully');
                    return true;
                }
            }
            return false;
        });
    };

    /**
     * Update PHP-specific elements to Go equivalents
     * Transforms PHP documentation links, button text, and icons to Go versions.
     * Replaces PHP icon with Go gopher SVG.
     * 
     * @returns {Promise<boolean>} Success status
     */
    GovelIgnitionUI.prototype.updatePhpToGo = async function() {
        return DOM.retry(() => {
            const phpDocsLink = DOM.select(CONFIG.SELECTORS.PHP_DOCS_LINK);
            if (!phpDocsLink) return false;

            // Update link URL
            phpDocsLink.href = CONFIG.URL_MAPPINGS['https://php.net/docs'];

            // Update button text
            const buttonText = phpDocsLink.querySelector('span:last-child');
            if (buttonText && buttonText.textContent === 'Docs') {
                buttonText.textContent = CONFIG.TEXT_REPLACEMENTS['Docs'];
            }

            // Replace PHP icon with Go gopher SVG
            var phpIcon = DOM.select(CONFIG.SELECTORS.PHP_ICON, phpDocsLink);
            if (phpIcon) {
                // Copy attributes from original SVG
                var newSvg = DOM.createSVGElement(CONFIG.GO_GOPHER_SVG);
                if (newSvg) {
                    // Preserve original SVG classes and attributes
                    Array.from(phpIcon.attributes).forEach(attr => {
                        if (attr.name !== 'data-icon') {
                            newSvg.setAttribute(attr.name, attr.value);
                        }
                    });
                    newSvg.setAttribute('data-icon', 'go');
                    // Set 14px size for gopher SVG
                    newSvg.setAttribute('width', '14');
                    newSvg.setAttribute('height', '14');
                    
                    phpIcon.parentNode.replaceChild(newSvg, phpIcon);
                    this.logger.info('PHP icon replaced with Go gopher SVG (14px)');
                }
            }

            this.logger.info('PHP to Go conversion completed');
            return true;
        });
    };

    /**
     * Update language version display from PHP to Go
     * Scans text content and replaces "PHP" with "Go" in version displays.
     * 
     * @returns {Promise<boolean>} Success status
     */
    GovelIgnitionUI.prototype.updateLanguageVersion = async function() {
        return DOM.retry(() => {
            const languageSpans = DOM.selectAll(CONFIG.SELECTORS.LANGUAGE_SPANS);
            let updated = false;

            languageSpans.forEach(span => {
                if (span.textContent && span.textContent.includes('PHP')) {
                    const trackingElement = span.querySelector(CONFIG.SELECTORS.TRACKING_WIDER);
                    if (trackingElement && trackingElement.textContent === 'PHP') {
                        trackingElement.textContent = CONFIG.TEXT_REPLACEMENTS['PHP'];
                        updated = true;
                    }
                }
            });

            if (updated) {
                this.logger.info('Language version updated from PHP to Go');
            }
            return updated;
        });
    };

    /**
     * Update source and documentation links
     * Converts Laravel/PHP-specific links to Go equivalents including:
     * - Laravel Ignition repository → GoVel repository
     * - Laravel website → Go website
     * - Flare documentation → Go documentation
     * 
     * @returns {Promise<boolean>} Success status
     */
    GovelIgnitionUI.prototype.updateSourceLinks = async function() {
        return DOM.retry(() => {
            let updated = false;

            // Update source repository link
            const sourceLink = DOM.select(CONFIG.SELECTORS.SOURCE_LINK);
            if (sourceLink) {
                sourceLink.href = CONFIG.URL_MAPPINGS['https://github.com/spatie/laravel-ignition'];
                updated = true;
                this.logger.info('Source link updated to GoVel repository');
            }

            // Update Laravel website link
            const laravelLink = DOM.select(CONFIG.SELECTORS.LARAVEL_LINK);
            if (laravelLink && laravelLink.textContent === 'Laravel') {
                laravelLink.href = CONFIG.URL_MAPPINGS['https://laravel.com'];
                laravelLink.textContent = CONFIG.TEXT_REPLACEMENTS['Laravel'];
                updated = true;
                this.logger.info('Laravel link updated to Go website');
            }

            // Update documentation links
            const docsLink = DOM.select(CONFIG.SELECTORS.DOCS_LINK);
            if (docsLink) {
                docsLink.href = CONFIG.URL_MAPPINGS['https://php.net/docs'];
                updated = true;
                this.logger.info('Documentation link updated to Go docs');
            }

            return updated;
        });
    };

    /**
     * Add Go gopher branding to Ignition logo area
     * Finds the Ignition logo and adds a subtle "Go" label to indicate
     * this is the Go version of Ignition.
     * 
     * @returns {Promise<boolean>} Success status
     */
    GovelIgnitionUI.prototype.addGoGopherIcon = async function() {
        return DOM.retry(() => {
            const ignitionLogo = DOM.select(CONFIG.SELECTORS.IGNITION_LOGO);
            if (!ignitionLogo || !ignitionLogo.parentElement) return false;

            // Check if Go label already exists
            const existingLabel = ignitionLogo.parentElement.querySelector('[data-govel-label="go"]');
            if (existingLabel) return true;

            // Create and style Go label
            const goLabel = document.createElement('span');
            goLabel.textContent = 'Go';
            goLabel.setAttribute('data-govel-label', 'go');
            
            // Apply styling
            Object.assign(goLabel.style, {
                fontSize: '10px',
                opacity: '0.6',
                marginLeft: '4px',
                fontWeight: 'bold',
                textTransform: 'uppercase',
                letterSpacing: '0.05em',
                color: 'currentColor',
                userSelect: 'none'
            });

            ignitionLogo.parentElement.appendChild(goLabel);
            this.logger.info('Go label added to Ignition logo');
            return true;
        });
    };

    /**
     * Get status of all modifications
     * @returns {Object} Status object with modification results
     */
    GovelIgnitionUI.prototype.getStatus = function() {
        return {
            initialized: this.initialized,
            modifications: Array.from(this.modifications),
            timestamp: new Date().toISOString()
        };
    };

    /**
     * Get logger instance for static methods
     * @returns {Logger} Logger instance
     */
    GovelIgnitionUI.getInstance = function() {
        return Logger.instance || new Logger();
    };
    
    // Also export to global scope for compatibility
    if (typeof window !== 'undefined') {
        window.GovelIgnitionUI = GovelIgnitionUI;
    }
    
    // Return the constructor for AMD
    return GovelIgnitionUI;
});
