# log

By creating logger interface, you can replace to any logging libraries after you implemented log.Infof(...) and so on.

Place log package to internal so that this package could not access from other libraries.
