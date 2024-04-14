import json
from models import MulipleChoice
from pydantic import ValidationError

TEMPLATE = """[INST] {{ .System }}
## Task
User:
{{ .Prompt }}
Output:
[/INST] """

SYSTEM = """You generate quizzes in hungarian from user supplied material. Output json only! Make sure 4 options are provided, and the solution can be any one or more of A, B, C, D.
## Example
User:
Magyarország állam Közép-Európában, a Kárpát-medence közepén. 1989 óta parlamentáris köztársaság. Északról Szlovákia, északkeletről Ukrajna, keletről és délkeletről Románia, délről Szerbia, délnyugatról Horvátország és Szlovénia, nyugatról pedig Ausztria határolja.
Output:
{"Question": "Melyik ország Magyarország északi szomszédja?", "Answers": ["Szlovénia", "Szlovákia", "Ukrajna", "Ausztria"], "Solutions": ["B"]}
"""


def try_parse_multiple_choice(data: str) -> MulipleChoice | None:
    try:
        response = json.loads(data)
        return MulipleChoice(
            question=response['Question'],
            options=response['Answers'],
            correct_options=response['Solutions'],
        )
    except json.JSONDecodeError or ValidationError or KeyError:
        return None
