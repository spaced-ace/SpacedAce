import langdetect as ld
from langdetect import DetectorFactory

DetectorFactory.seed = 0

ALLOWED_LANGUAGES = ['hu', 'en']

FALLBACK_LANG = 'en'


def detect_language(text: str) -> str:
    lang = ld.detect_langs(text)
    langs = [l.lang for l in lang]
    for l in langs:
        if l in ALLOWED_LANGUAGES:
            return l
    return FALLBACK_LANG
