/**
 * DOM utility functions
 * 
 * Provides safe DOM operations, element querying, and utility functions
 * for the GoVel Ignition UI system. Includes retry mechanisms for 
 * handling asynchronous DOM operations.
 * 
 * @version 1.0.0
 * @author GoVel Team
 * @created 2025-09-11
 * @requires RequireJS
 * @module dom
 */

define([], function() {
    'use strict';
    
    /**
     * DOM utility functions
     * @namespace DOM
     */
    var DOM = {
        /**
         * Safe element selector with error handling
         * @param {string} selector - CSS selector
         * @param {Element} [context=document] - Context element
         * @returns {Element|null} Found element or null
         */
        select: function(selector, context) {
            context = context || document;
            try {
                return context.querySelector(selector);
            } catch (error) {
                console.error('[GoVel Ignition] Invalid selector: ' + selector, error);
                return null;
            }
        },

        /**
         * Safe element selector for multiple elements
         * @param {string} selector - CSS selector
         * @param {Element} [context=document] - Context element
         * @returns {NodeList} Found elements
         */
        selectAll: function(selector, context) {
            context = context || document;
            try {
                return context.querySelectorAll(selector);
            } catch (error) {
                console.error('[GoVel Ignition] Invalid selector: ' + selector, error);
                return [];
            }
        },

        /**
         * Check if element exists and is visible
         * @param {Element} element - Element to check
         * @returns {boolean} True if element exists and is visible
         */
        isVisible: function(element) {
            if (!element) return false;
            var style = window.getComputedStyle(element);
            return style.display !== 'none' && style.visibility !== 'hidden' && style.opacity !== '0';
        },

        /**
         * Create SVG element from string
         * @param {string} svgContent - SVG content string
         * @returns {SVGElement|null} Created SVG element
         */
        createSVGElement: function(svgContent) {
            try {
                var parser = new DOMParser();
                // Use 14x14 viewBox for Go gopher SVG, 24x24 for others
                var viewBox = svgContent.indexOf('GO_GOPHER') !== -1 || svgContent.indexOf('#8CC5E7') !== -1 ? '0 0 24 24' : '0 0 24 24';
                var svgDoc = parser.parseFromString('<svg viewBox="' + viewBox + '" width="14" height="14" xmlns="http://www.w3.org/2000/svg">' + svgContent + '</svg>', 'image/svg+xml');
                return svgDoc.documentElement;
            } catch (error) {
                console.error('[GoVel Ignition] Failed to create SVG element', error);
                return null;
            }
        },

        /**
         * Wait for DOM to be ready
         * @param {Function} callback - Function to execute when DOM is ready
         */
        ready: function(callback) {
            if (document.readyState !== 'loading') {
                callback();
            } else {
                document.addEventListener('DOMContentLoaded', callback);
            }
        },

        /**
         * Retry operation with exponential backoff
         * @param {Function} operation - Operation to retry
         * @param {number} [maxRetries=10] - Maximum retry attempts
         * @param {number} [delay=50] - Initial delay in ms
         * @returns {Promise} Promise that resolves when operation succeeds
         */
        retry: function(operation, maxRetries, delay) {
            maxRetries = maxRetries || 10;
            delay = delay || 50;
            
            return new Promise(function(resolve, reject) {
                function attemptOperation(attempt) {
                    try {
                        var result = operation();
                        if (result) {
                            resolve(result);
                            return;
                        }
                    } catch (error) {
                        console.warn('[GoVel Ignition] Retry ' + (attempt + 1) + '/' + maxRetries + ' failed:', error.message);
                    }
                    
                    if (attempt >= maxRetries - 1) {
                        reject(new Error('Operation failed after ' + maxRetries + ' retries'));
                        return;
                    }
                    
                    setTimeout(function() {
                        attemptOperation(attempt + 1);
                    }, delay * Math.pow(2, attempt));
                }
                
                attemptOperation(0);
            });
        }
    };
    
    // Return the DOM utilities object
    return DOM;
});
