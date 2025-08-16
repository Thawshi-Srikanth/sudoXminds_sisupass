import json


class PrettyJSONEncoder(json.JSONEncoder):
    def __init__(self, *args, **kwargs):
        # Provide defaults, but allow override via kwargs
        kwargs.setdefault("indent", 4)
        kwargs.setdefault("sort_keys", True)
        super().__init__(*args, **kwargs)
