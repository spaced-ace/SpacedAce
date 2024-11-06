import abc


class Provider(metaclass=abc.ABCMeta):
    """Groups common functions of LLM providers"""

    CHAT_COMPLETION_ENDPOINT: str

    @staticmethod
    @abc.abstractmethod
    def parse_chat_completion_response(response: dict) -> str:
        """Accesses the llm's text output from the response object"""


class Ollama(Provider):
    CHAT_COMPLETION_ENDPOINT = '/api/chat'

    def __init__(self) -> None:
        pass

    @staticmethod
    def parse_chat_completion_response(response: dict) -> str:
        return response['message']['content']


class OpenAI(Provider):
    CHAT_COMPLETION_ENDPOINT = '/chat/completions'

    def __init__(self) -> None:
        pass

    @staticmethod
    def parse_chat_completion_response(response: dict) -> str:
        return response['choices'][0]['message']['content']
