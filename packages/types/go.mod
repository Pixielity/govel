module govel/types

go 1.25

// Map the main module to the src directory
replace govel/types => ./src

replace (
    // Map the main module to the src directory
    govel/types => ./src

    // Inter-package dependencies
    govel/application => ../application
    govel/config => ../config
    govel/container => ../container
    govel/cookie => ../cookie
    govel/encryption => ../encryption
    govel/hashing => ../hashing
    govel/logger => ../logger
    govel/pipeline => ../pipeline
    govel/support => ../support

)
