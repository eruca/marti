    /**  Map is a general map object for storing key value pairs
     *  @param m - default set of properties
     */
var Map =function(m) {
    var map;
    if (typeof m == 'undefined')
         map = new Array();
    else
         map = m;
    /**
     * Get a list of the keys to check
     */
    this.keys = function() {
        var _keys = new Array();
        for (var _i in map){
            _keys.push(_i);
        }
        return _keys;
    };
    /**
    * Put stores the value in the table
    * @param key the index in the table where the value will be stored
    * @param value the value to be stored
    */
    this.put = function(key,value) {
        map[key] = value;
    };
    /**
    * Return the value stored in the table
    * @param key the index of the value to retrieve
    */
    this.get = function(key) {
        return map[key];
    };

    /**
     * Remove the value from the table
     * @param key the index of the value to be removed
     */
    this.remove =  function(key) {
        map[key]=null;
        delete map[key];
    };
    /**
    *  Clear the table
    */
    this.clear = function() {
        delete map;
        map = new Array();
    };
};

//example
// var m=new Map();
// m.put("id","1000");
// m.put("name","张三");