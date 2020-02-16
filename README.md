# Go-socket.io-parser

It is copied repository from - [url](https://github.com/googollee/go-engine.io)

# Documentation

Official documentation engine.io protocol - [url](https://github.com/socketio/engine.io-protocol)

# Gratitude

More thanks by [@googolle](https://www.github.com/googolle) for created and migrated engine.io-parser protocol to pure go

# Benchmarks

## Machine properties:
 
Model Name: MacBook Pro
Model Identifier: MacBookPro14,1
Processor Name: Dual-Core Intel Core i7
Processor Speed: 2,5 GHz
Number of Processors: 1
Total Number of Cores: 2
L2 Cache (per Core): 256 KB
L3 Cache: 4 MB
Hyper-Threading Technology: Enabled
Memory: 16 GB


[Node.js]()
--------
encode packet as string x 96,175 ops/sec ±58.02% (59 runs sampled)
encode packet as binary x 98,193 ops/sec ±9.07% (51 runs sampled)
encode payload as string x 68,634 ops/sec ±9.78% (48 runs sampled)
encode payload as binary x 47,832 ops/sec ±12.54% (47 runs sampled)
decode packet from string x 26,802,947 ops/sec ±5.60% (74 runs sampled)
decode packet from binary x 5,813,963 ops/sec ±2.63% (84 runs sampled)
decode payload from string x 81,404 ops/sec ±11.06% (34 runs sampled)
decode payload from binary x 62,011 ops/sec ±36.14% (40 runs sampled)

[Golang]()
