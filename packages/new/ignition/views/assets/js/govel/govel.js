/**
 * GoVel Ignition UI Entry Point
 * 
 * Main entry point for the GoVel Ignition UI customization system.
 * Loads and initializes all required modules and components.
 * 
 * @version 1.0.0
 * @author GoVel Team
 * @created 2025-09-11
 * @requires RequireJS, CONFIG, Logger, DOM, GovelIgnitionUI
 */

define(['./govel-ui'], function(GovelIgnitionUI) {
    'use strict';

    /**
     * Load and initialize the GoVel Ignition UI system
     * @private
     */
    async function initialize() {
        try {
            // Create and initialize the UI customizer
            var uiCustomizer = new GovelIgnitionUI({
                debug: (window.location && window.location.search && window.location.search.includes('debug=1')) || false
            });

            // Initialize the UI customization
            await uiCustomizer.init();

            // Expose to global scope for debugging if debug mode is enabled
            if (window.location && window.location.search && window.location.search.includes('debug=1')) {
                window.GovelIgnitionDebug = uiCustomizer;
                console.log('[GoVel Ignition] Debug mode enabled. UI customizer exposed as window.GovelIgnitionDebug');
            }

        } catch (error) {
            console.error('[GoVel Ignition] Failed to initialize UI system:', error);
        }
    }

    // Initialize when DOM is ready or immediately if already loaded
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initialize);
    } else {
        // Small delay to ensure all scripts are loaded
        setTimeout(initialize, 10);
    }
    
    // Return initialization function for AMD
    return {
        init: initialize
    };
});
