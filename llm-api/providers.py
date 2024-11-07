import abc
import httpx


class Provider(metaclass=abc.ABCMeta):
    """Groups common functions of LLM providers"""

    NAME: str

    @abc.abstractmethod
    async def get_model_response(self, messages: list[dict]) -> str:
        """Returns the LLM response for an openai style conversation"""


class Ollama(Provider):
    NAME = 'ollama'

    def __init__(
        self, client: httpx.AsyncClient, base_url: str, model_name: str
    ) -> None:
        self.client = client
        self.model_name = model_name
        self.base_url = base_url
        self.url = f'{base_url}/api/chat'

    async def get_model_response(self, messages: list[dict]) -> str:
        """Returns the LLM response for an openai style conversation"""
        payload = {
            'temperature': 0.4,
            'model': self.model_name,
            'stream': False,
            'messages': messages,
            'max_new_tokens': 500,
        }
        res = await self.client.post(
            self.url,
            json=payload,
        )
        res.raise_for_status()
        return res.json()['message']['content']


class OpenAI(Provider):
    NAME = 'openai'

    def __init__(
        self,
        client: httpx.AsyncClient,
        base_url: str,
        api_key: str,
        model_name: str,
    ) -> None:
        self.client = client
        self.api_key = api_key
        self.model_name = model_name
        self.base_url = base_url
        self.url = f'{base_url}/chat/completions'

    async def get_model_response(self, messages: list[dict]) -> str:
        """Returns the LLM response for an openai style conversation"""
        headers = {'Authorization': f'Bearer {self.api_key}'}
        payload = {
            'temperature': 0.4,
            'model': self.model_name,
            'stream': False,
            'messages': messages,
        }
        res = await self.client.post(
            self.url,
            headers=headers,
            json=payload,
        )
        res.raise_for_status()
        return res.json()['choices'][0]['message']['content']


class Google(Provider):
    NAME = 'google'
    _HARM_CATEGORIES = [
        'HARM_CATEGORY_HATE_SPEECH',
        'HARM_CATEGORY_SEXUALLY_EXPLICIT',
        'HARM_CATEGORY_DANGEROUS_CONTENT',
        'HARM_CATEGORY_HARASSMENT',
    ]

    def __init__(
        self,
        client: httpx.AsyncClient,
        api_key: str,
        model_name: str,
    ) -> None:
        self.client = client
        self.api_key = api_key
        self.model_name = model_name
        self.url = f'https://generativelanguage.googleapis.com/v1beta/models/{self.model_name}:generateContent?key={self.api_key}'

    async def get_model_response(self, messages: list[dict]) -> str:
        """Returns the LLM response for an openai style conversation"""
        system, convo = Google._openai_conversation_to_google(messages)
        gen_conf = {
            'temperature': 0.4,
        }
        headers = {'Content-Type': 'application/json'}
        payload = {
            'safetySettings': [
                {
                    'category': cat,
                    'threshold': 'BLOCK_NONE',
                }
                for cat in Google._HARM_CATEGORIES
            ],
            'generationConfig': gen_conf,
            'contents': convo,
        }
        if system is not None:
            payload['system_instruction'] = system

        response = await self.client.post(
            self.url,
            json=payload,
            headers=headers,
        )
        response.raise_for_status()
        return response.json()['candidates'][0]['content']['parts'][0]['text']

    @staticmethod
    def _openai_conversation_to_google(
        convo: list[dict],
    ) -> tuple[dict | None, list[dict]]:
        if len(convo) == 0:
            return None, []
        convo_cpy = convo.copy()
        system = None
        if convo[0]['role'] == 'system':
            system = {'parts': {'text': convo_cpy[0]['content']}}
            convo_cpy = convo_cpy[1:]
        convo_transformed = [
            {
                'role': 'model' if msg['role'] == 'assistant' else 'user',
                'parts': [{'text': msg['content']}],
            }
            for msg in convo_cpy
        ]
        return system, convo_transformed
