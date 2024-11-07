import json
from models import MulipleChoice, SingleChoice, TrueOrFalse
from pydantic import ValidationError

import models

SYSTEM_EN = 'You are a helpful assistant to a teacher, who creates test questions for students in json format.'

PROMPT_EN = """
Create a {} question based on the context.
Example:
<context>{}</context>
<output>{}</output>
Task:
<context>{}</context>
"""
PROMPT_HU = """
Írj egy {} kérdést a kontextus alapján.
Példa:
<context>{}</context>
<output>{}</output>
Feladat:
<context>{}</context>
"""

RESPONSE = """<output>{}</output>"""

EXAMPLE_CONTEXT_EN = """
The Nobel Prize in Literature (here meaning for literature; Swedish: Nobelpriset i litteratur) is a Swedish literature prize that is awarded annually, since 1901, to an author from any country who has, in the words of the will of Swedish industrialist Alfred Nobel, "in the field of literature, produced the most outstanding work in an idealistic direction"
"""
MULTI_EXAMPLE_EN = """
{
    "question": "Which of the following statements are true about the Nobel Prize in Literature?",
    "answers": [
        "It is awarded annually.",
        "It is only awarded to Swedish authors.",
        "It has been awarded since 1901.",
        "It is given for outstanding work in the field of literature."
        ],
    "solutions": ["A", "C", "D"]
}"""
SINGLE_EXAMPLE_EN = """
{
    "question": "Which of the following statements are true about the Nobel Prize in Literature?",
    "answers": [
        "It is awarded to an author from Sweden.",
        "It is awarded to an author from any country for producing outstanding work in an idealistic direction.",
        "It is awarded to the best-selling author of the year.",
        "It is awarded to an author for writing about Swedish history."
    ],
    "solution": "B"
}"""
BOOLEAN_EXAMPLE_EN = """
{
    "question":"The Nobel Prize in Literature is awarded annually to authors only from Sweden.",
    "solution":false
}"""

EXAMPLE_CONTEXT_HU = """Nobel-díjat a svéd kémikus és feltaláló Alfred Nobel alapította. Nobel 1895 november 27-én kelt végrendeletében rendelkezett úgy, hogy vagyonának kamataiból évről évre részesedjenek a fizika, kémia, fiziológia és orvostudomány, továbbá az irodalom legjobbjai és az a személy, aki a békéért tett erőfeszítéseivel a díjat kiérdemli."""

BOOLEAN_EXAMPLE_HU = """
{
    "question":"Nobel-díjat csak a svéd kémikusok és feltalálók kaphatnak meg.",
    "solution":false
}"""

SINGLE_EXAMPLE_HU = """
{
    "question": "Mi volt Alfred Nobel végrendeletének célja a Nobel-díjjal kapcsolatban?",
    "answers": [
        "Csak svéd tudósoknak adják át.",
        "A fizika, kémia, fiziológia, orvostudomány, irodalom legjobbjait és a békéért küzdő személyt jutalmazzák.",
        "Csak irodalmi teljesítményért ítélik oda.",
        "A legújabb találmányokat jutalmazzák."
    ],
    "solution": "B"
}"""

MULTI_EXAMPLE_HU = """
{
    "question": "Mely állítások igazak a Nobel-díjjal kapcsolatban?",
    "answers": [
    "Alfred Nobel alapította a díjat.",
    "A díjat csak fizikai teljesítményért ítélik oda.",
    "A végrendeletében rendelkezett a díj alapításáról.",
    "A békéért tett erőfeszítéseket is jutalmazzák."
        ],
    "solution": ["A", "C", "D"]
}"""

SYSTEM_HU = 'Segítőkész asszisztens vagy egy tanárnak, aki tesztkérdéseket készít a diákok számára json formátumban.'


def format_question(context: str, q_type: str, lang: str) -> list[dict]:
    """Formats a question into a list of conversation turns"""

    if lang == 'en':
        question_type = (
            'boolean'
            if q_type == models.TRUE_OR_FALSE
            else (
                'multiple choice single answer (4 options)'
                if q_type == models.SINGLE_CHOICE
                else 'multiple choice multiple answers (4 options)'
            )
        )
    elif lang == 'hu':
        question_type = (
            'igaz/hamis'
            if q_type == models.TRUE_OR_FALSE
            else (
                'egy válaszlehetőséges (4 opciós)'
                if q_type == models.SINGLE_CHOICE
                else 'több válaszlehetőséges (4 opciós)'
            )
        )
    else:
        raise ValueError('Unsupported language')

    if lang == 'en':
        example = (
            BOOLEAN_EXAMPLE_EN
            if q_type == models.TRUE_OR_FALSE
            else SINGLE_EXAMPLE_EN
            if q_type == models.SINGLE_CHOICE
            else MULTI_EXAMPLE_EN
        )
        prompt = PROMPT_EN.format(
            question_type,
            EXAMPLE_CONTEXT_EN,
            example,
            context,
        )
    elif lang == 'hu':
        example = (
            BOOLEAN_EXAMPLE_HU
            if q_type == models.TRUE_OR_FALSE
            else SINGLE_EXAMPLE_HU
            if q_type == models.SINGLE_CHOICE
            else MULTI_EXAMPLE_HU
        )
        prompt = PROMPT_HU.format(
            question_type,
            EXAMPLE_CONTEXT_HU,
            example,
            context,
        )
    convo = [
        {
            'role': 'system',
            'content': SYSTEM_EN if lang == 'en' else SYSTEM_HU,
        },
        {'role': 'user', 'content': prompt},
    ]
    return convo


def strip_response(response: str) -> str:
    s = response.split('<output>')[-1]
    s = s.split('</output>')[0]
    s = s.strip()
    return s


def try_parse_multiple_choice(data: str) -> MulipleChoice | None:
    stripped = strip_response(data)
    try:
        response = json.loads(stripped)
        response = {k.lower(): v for k, v in response.items()}
        return MulipleChoice(
            question=response['question'],
            options=response['answers'],
            correct_options=response['solution'],
        )
    except json.JSONDecodeError or ValidationError or KeyError:
        print(stripped)
        return None


def try_parse_single_choice(data: str) -> SingleChoice | None:
    stripped = strip_response(data)
    try:
        response = json.loads(stripped)
        response = {k.lower(): v for k, v in response.items()}
        return SingleChoice(
            question=response['question'],
            options=response['answers'],
            correct_option=response['solution'],
        )
    except json.JSONDecodeError or ValidationError or KeyError:
        print(stripped)
        return None


def try_parse_true_or_false(data: str) -> TrueOrFalse | None:
    stripped = strip_response(data)
    try:
        response = json.loads(stripped)
        response = {k.lower(): v for k, v in response.items()}
        return TrueOrFalse(
            question=response['question'],
            correct_option=response['solution'],
        )
    except json.JSONDecodeError or ValidationError or KeyError:
        print(stripped)
        return None
