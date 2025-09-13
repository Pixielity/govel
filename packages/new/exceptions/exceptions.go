// Package exceptions provides a comprehensive exception handling system for GoVel applications.
// This package follows Laravel's exception pattern with ISP-compliant interfaces,
// providing HTTP-aware exceptions with centralized handling, flexible rendering,
// and comprehensive solution support.
package exceptions

// Re-export core interfaces for backward compatibility and ease of use
import (
	"govel/exceptions/core"
	"govel/exceptions/core/solution"
	"govel/exceptions/helpers"
	httpExceptions "govel/exceptions/http"
	"govel/exceptions/interfaces"
	solutionInterface "govel/exceptions/interfaces/solution"
	httpSolutions "govel/exceptions/solutions/http"
	"govel/exceptions/solutions/runnable"
	"govel/exceptions/solutions/providers"
)

// =============================================================================
// Type Aliases for Backward Compatibility
// =============================================================================

// ExceptionInterface is the main interface for all GoVel exceptions
type ExceptionInterface = interfaces.ExceptionInterface

// HTTPable interface for HTTP-related functionality
type HTTPable = interfaces.HTTPable

// Contextable interface for context management
type Contextable = interfaces.Contextable

// Renderable interface for rendering exceptions
type Renderable = interfaces.Renderable

// Stackable interface for stack trace functionality
type Stackable = interfaces.Stackable

// Solutionable interface for solution-related functionality
type Solutionable = interfaces.Solutionable

// Solution interfaces
type Solution = solutionInterface.Solution
type ProvidesSolution = solutionInterface.ProvidesSolution
type RunnableSolution = solutionInterface.RunnableSolution
type HasSolutionsForThrowable = solutionInterface.HasSolutionsForThrowable
type SolutionProvider = solutionInterface.SolutionProvider

// Core types
type Exception = core.Exception
type BaseSolution = solution.BaseSolution
type SolutionProviderRepository = solution.SolutionProviderRepository

// HTTP Exception types
type BadRequestException = httpExceptions.BadRequestException
type UnauthorizedException = httpExceptions.UnauthorizedException
type ForbiddenException = httpExceptions.ForbiddenException
type NotFoundException = httpExceptions.NotFoundException
type MethodNotAllowedException = httpExceptions.MethodNotAllowedException
type ConflictException = httpExceptions.ConflictException
type UnprocessableEntityException = httpExceptions.UnprocessableEntityException
type TooManyRequestsException = httpExceptions.TooManyRequestsException
type InternalServerErrorException = httpExceptions.InternalServerErrorException
type ServiceUnavailableException = httpExceptions.ServiceUnavailableException

// =============================================================================
// Core Functions
// =============================================================================

// NewException creates a new base exception
var NewException = core.NewException

// NewBaseSolution creates a new base solution
var NewBaseSolution = solution.NewBaseSolution

// NewSolutionProviderRepository creates a new solution provider repository
var NewSolutionProviderRepository = solution.NewSolutionProviderRepository

// =============================================================================
// HTTP Exception Constructors
// =============================================================================

// NewBadRequestException creates a new 400 Bad Request exception
var NewBadRequestException = httpExceptions.NewBadRequestException

// NewUnauthorizedException creates a new 401 Unauthorized exception
var NewUnauthorizedException = httpExceptions.NewUnauthorizedException

// NewForbiddenException creates a new 403 Forbidden exception
var NewForbiddenException = httpExceptions.NewForbiddenException

// NewNotFoundException creates a new 404 Not Found exception
var NewNotFoundException = httpExceptions.NewNotFoundException

// NewMethodNotAllowedException creates a new 405 Method Not Allowed exception
var NewMethodNotAllowedException = httpExceptions.NewMethodNotAllowedException

// NewConflictException creates a new 409 Conflict exception
var NewConflictException = httpExceptions.NewConflictException

// NewUnprocessableEntityException creates a new 422 Unprocessable Entity exception
var NewUnprocessableEntityException = httpExceptions.NewUnprocessableEntityException

// NewTooManyRequestsException creates a new 429 Too Many Requests exception
var NewTooManyRequestsException = httpExceptions.NewTooManyRequestsException

// NewInternalServerErrorException creates a new 500 Internal Server Error exception
var NewInternalServerErrorException = httpExceptions.NewInternalServerErrorException

// NewServiceUnavailableException creates a new 503 Service Unavailable exception
var NewServiceUnavailableException = httpExceptions.NewServiceUnavailableException

// =============================================================================
// Solution Constructors
// =============================================================================

// HTTP Solutions
var NewBadRequestSolution = httpSolutions.NewBadRequestSolution
var NewUnauthorizedSolution = httpSolutions.NewUnauthorizedSolution
var NewForbiddenSolution = httpSolutions.NewForbiddenSolution
var NewNotFoundSolution = httpSolutions.NewNotFoundSolution
var NewMethodNotAllowedSolution = httpSolutions.NewMethodNotAllowedSolution
var NewValidationErrorSolution = httpSolutions.NewValidationErrorSolution
var NewTooManyRequestsSolution = httpSolutions.NewTooManyRequestsSolution
var NewInternalServerErrorSolution = httpSolutions.NewInternalServerErrorSolution
var NewServiceUnavailableSolution = httpSolutions.NewServiceUnavailableSolution
var NewConflictSolution = httpSolutions.NewConflictSolution

// Runnable Solutions
var NewGenerateAppKeySolution = runnable.NewGenerateAppKeySolution
var NewCreateDirectorySolution = runnable.NewCreateDirectorySolution
var NewInstallDependencySolution = runnable.NewInstallDependencySolution
var NewFixPermissionsSolution = runnable.NewFixPermissionsSolution

// Solution Providers
var NewHTTPExceptionProvider = providers.NewHTTPExceptionProvider
var NewCommonRunnableSolutionsProvider = providers.NewCommonRunnableSolutionsProvider

// =============================================================================
// Helper Functions
// =============================================================================

// Abort creates and returns a new exception with the given status code and message
var Abort = helpers.Abort

// AbortIf creates and returns a new exception if the given condition is true
var AbortIf = helpers.AbortIf

// AbortUnless creates and returns a new exception if the given condition is false
var AbortUnless = helpers.AbortUnless

// Shortcut functions
var Abort400 = helpers.Abort400
var Abort401 = helpers.Abort401
var Abort403 = helpers.Abort403
var Abort404 = helpers.Abort404
var Abort500 = helpers.Abort500
