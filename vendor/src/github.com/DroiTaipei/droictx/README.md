# droictx

## Abstract 
* The Droi Context package, it supply the Context Type and utilities.

## Usage
* `GetContextFromPeeker`
  * It designed for getting Context From a Peeker ( with Peek method) 
  * Usually, a peeker is a *fasthttp.RequestHeader
  * It will convert I-Fields, e.g. "X-Droi-AppID" -> "Aid" 
  * Sample Code 

```
func GetContext(r *fasthttp.RequestHeader) Context {
	return GetContextFromPeeker(r)
}
```

* `SetHTTPHeaders`
  * It designed for set HTTP Headers From Context with a setter ( with Set method) 
  * Usually, a setter is a *fasthttp.RequestHeader
  * Sample Code 

```
var c droictx.Context
var r r *fasthttp.RequestHeader
c.SetHeaderFromContext(r)
# Header should like X-Droi-AppID: xxxxxx
```



* `(c *Context) HeaderMap()`
  * It get the I-Fields values from Context, return a map with HTTP Header Key
  * e.g.
  * Stored in Context
  
 ```
 {
   "Aid": "asdqwezxc123"
 }
 
 ```
 
  * Return amp
  

 ```
 {
   "X-Droi-AppID" : "asdqwezxc123"
 }
 
 ```
  


