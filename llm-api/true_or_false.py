import json
from models import TrueOrFalse
from pydantic import ValidationError

TEMPLATE = """[INST] {{ .System }}
## Task
User:
{{ .Prompt }}
Output:
[/INST] """

SYSTEM = """You generate quizzes in hungarian from user supplied material. Output json only! Make sure the solution is either true or false.
## Example
User:
Magyarország állam Közép-Európában, a Kárpát-medence közepén. 1989 óta parlamentáris köztársaság. Északról Szlovákia, északkeletről Ukrajna, keletről és délkeletről Románia, délről Szerbia, délnyugatról Horvátország és Szlovénia, nyugatról pedig Ausztria határolja.
Output:
{"Question": "Ukrajna Magyarország északi szomszédja. Igaz vagy hamis?", "Solution": false}
"""


def try_parse_true_or_false(data: str) -> TrueOrFalse | None:
    try:
        response = json.loads(data)
        return TrueOrFalse(
            question=response['Question'],
            correct_option=response['Solution'],
        )
    except json.JSONDecodeError or ValidationError or KeyError:
        return None
