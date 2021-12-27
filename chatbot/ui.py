try:
    # Python 3
    from http.server import HTTPServer, SimpleHTTPRequestHandler, test as test_orig
    import sys
    def test (*args):
        test_orig(*args, port=int(sys.argv[1]) if len(sys.argv) > 1 else 8000)
except ImportError: # Python 2
    from BaseHTTPServer import HTTPServer, test
    from SimpleHTTPServer import SimpleHTTPRequestHandler
 
class CORSRequestHandler (SimpleHTTPRequestHandler):
    def end_headers (self):
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header("Access-Control-Allow-Headers", "X-Requested-With");
        self.send_header("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS");
        self.send_header("X-Powered-By", ' 3.2.1')
        # self.send_header("Content-Type", "application/json;charset=utf-8");
        SimpleHTTPRequestHandler.end_headers(self)
 
if __name__ == '__main__':
    test(CORSRequestHandler, HTTPServer)