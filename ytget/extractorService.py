#!/usr/bin/env python3

from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse, parse_qs
from youtube_dl import YoutubeDL as YT

class HTTPServer_RequestHandler(BaseHTTPRequestHandler):

    def do_GET(self):

        # get parameters
        query_components = parse_qs(urlparse(self.path).query)
        url = str(query_components.get("input")[0])
        if url == None:
            self.bad_request("missing 'input' parameter")
            return

        # run YoutubeDL url extractor
        print(url)
        out = self.url_extractor(url)
        if out == None:
            self.bad_request("youtube_dl extractor failed")
            return

        # generate response
        self.send_response(200)
        self.send_header('Content-type','text/html')
        self.end_headers()

        # Send message back to client
        # Write content as utf-8 data
        self.wfile.write(bytes(out, "utf8"))
        return

    def url_extractor(self, url):
        yt = YT({ 'quiet': True, 'skip_download': True })
        try:
            info_dict = yt.extract_info(url)
            for format in info_dict['formats']:
                if format['vcodec'] == "none":
                    result = format['url']
                    break
            return result
        except:
            return None

    def bad_request(self, message):
        self.send_response(400)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(bytes(message, "utf8"))

if __name__ == "__main__":
    server_address = ("localhost", 8081)
    httpd = HTTPServer(server_address, HTTPServer_RequestHandler)
    print('running server...')
    httpd.serve_forever()
