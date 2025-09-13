/**
 * Logger utility for development and debugging
 * 
 * Provides centralized logging functionality for the GoVel Ignition UI system.
 * Supports debug mode toggling and consistent log formatting across the application.
 * 
 * @version 1.0.0
 * @author GoVel Team
 * @created 2025-09-11
 * @requires RequireJS
 * @module logger
 */

define([], function() {
    'use strict';
    
    /**
     * Logger utility for development and debugging
     * @class Logger
     */
    function Logger(isDebugMode) {
        this.debug = isDebugMode || false;
        this.prefix = '[GoVel Ignition]';
    }

    /**
     * Log info message
     * @param {string} message - The message to log
     * @param {...any} args - Additional arguments
     */
    Logger.prototype.info = function(message) {
        if (this.debug) {
            var args = Array.prototype.slice.call(arguments, 1);
            console.log.apply(console, [this.prefix + ' ' + message].concat(args));
        }
    };

    /**
     * Log warning message
     * @param {string} message - The message to log
     * @param {...any} args - Additional arguments
     */
    Logger.prototype.warn = function(message) {
        var args = Array.prototype.slice.call(arguments, 1);
        console.warn.apply(console, [this.prefix + ' ' + message].concat(args));
    };

    /**
     * Log error message
     * @param {string} message - The message to log
     * @param {...any} args - Additional arguments
     */
    Logger.prototype.error = function(message) {
        var args = Array.prototype.slice.call(arguments, 1);
        console.error.apply(console, [this.prefix + ' ' + message].concat(args));
    };
    
    // Return the Logger constructor
    return Logger;
});
