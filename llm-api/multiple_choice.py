import json
from models import MulipleChoice
from pydantic import ValidationError


def try_parse_multiple_choice(data: str) -> MulipleChoice | None:
    try:
        response = json.loads(data)
        return MulipleChoice(
            question=response['Question'],
            options=response['Answers'],
            correct_option=response['Solution'],
        )
    except json.JSONDecodeError or ValidationError or KeyError:
        return None
