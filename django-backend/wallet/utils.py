from functools import wraps
from graphql import GraphQLError

def login_required(func):
    @wraps(func)
    def wrapper(root, info, *args, **kwargs):
        user = info.context.user
        if user.is_anonymous:
            raise GraphQLError("Authentication required")
        return func(root, info, *args, **kwargs)
    return wrapper
