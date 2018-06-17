import sys
from django.core.validators import URLValidator
from youtube_dl import YoutubeDL as YT

yt = YT({ 'quiet': True, 'skip_download': True })

def get_url(url):

    try:
        info_dict = yt.extract_info(url)
        for format in info_dict['formats']:
            if format['vcodec'] == "none":
                result = format['url']
                break

        return result

    except:
        return "youtube_dl extractor failed"

target_url = get_url(sys.argv[1])
validator = URLValidator()

try:
    validator(target_url)
    print(target_url)
except:
    exit()
