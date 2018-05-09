package Status

const(
/// 2xx Success
  OK  int = 200
  Created int= 201
  Accepted int= 202
  NoContent int= 204
  ResetContent int= 205

/// 3xx Redirection
  MultipleChoices int= 300
  MovedPermanently int= 301
  Found int= 302
  SeeOther int= 303
  NotModified int= 304

/// 4xx Client Error
  BadRequest int= 400
  Unauthorized int= 401
  Forbidden int= 403
  NotFound int= 404
  MethodNotAllowed int= 405
  NotAcceptable int= 406
  ProxyAuthenticationRequired int= 407
  RequestTimeout int= 408
  Conflict int= 409
  Gone int= 410
  LengthRequired int= 411
  PreconditionFailed int= 412
  Locked int= 423
  PreconditionRequired int= 428
  TooManyRequests int= 429

/// 4xx Client Error - CUSTOM Game Errors
  ServiceNotFound int= 469
  ConnectionNotReady int= 470
  RequestNotStarted int= 471
  ConnectionError int= 472
  InvalidResponseHandle int= 473
  ResponseNotReady int= 474
  GetDataFailed int= 475
  ConnectionCanceled int= 476

  EncodeDataError int= 480
  DecodeDataError int= 481
  EngineImplementedYet int= 490
  RequestValidationError int= 491

/// 5xx Server Error
  InternalServerError int= 500
  NotImplemented int= 501
  BadGateway int= 502
  ServiceUnavailable int= 503
  GatewayTimeout int= 504

/// ETSv2 Errors
/// https://phabricator.gameloft.org/w/etsv2/etsv2-error-codes/
  GGIbanned_RejectedByAdmin int= 990
  GGIbanned_BlockedStopSending int= 991
  GGIbanned_BlockedDelete int= 992
  GGIbanned_BlockedDisableTracking int= 993
  ServerIsShuttingDown int= 994
  BatchStorageRefused int= 995
  RejectedEvents int= 996
  RejectedPackage int= 997
  RejectedHeaders int= 998
  TimerExpired int= 999
)