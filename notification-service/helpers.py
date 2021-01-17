def get_params_if_in_query(request, params):
    _dict = {}
    for param in params:
        value = request.args.get(param)
        if value:
            _dict[param] = value
    return _dict


def get_params_if_in_form(request, params):
    _dict = {}
    for param in params:
        try:
            value = request.form[param]
            _dict[param] = value
        except KeyError:
            continue
    return _dict
