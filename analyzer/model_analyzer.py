import argparse
import json
import os
import sys
import typing
from dataclasses import asdict, dataclass, is_dataclass
from pathlib import Path

from django.apps.registry import apps
from django.conf import ENVIRONMENT_VARIABLE, LazySettings
from django.db.models.query_utils import DeferredAttribute


class DCJSONEncoder(json.JSONEncoder):
    def default(self, o):
        if is_dataclass(o):
            return asdict(o)
        return super().default(o)


@dataclass
class FieldLookup:
    name: str
    type: str


@dataclass
class FieldType:
    name: str
    lookups: list[FieldLookup]


@dataclass
class ModelField:
    name: str
    type_name: str
    opts: dict[str, typing.Any]


@dataclass
class DjangoModel:
    app: str
    name: str
    fields: list[ModelField]


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Django project analyzer for django-ls")
    parser.add_argument("workspace_path", type=Path)
    parser.add_argument("-s", "--settings", type=str, help="Default settings module, overriden by env var.")
    parser.add_argument("-o", "--output", type=Path, help="Path to output file, uses stdout if not specified.")
    return parser.parse_args()


def run_analyzer(workspace_path: Path, settings_module: str | None = None, output: Path | None = None):
    sys.path.append(str(workspace_path))
    if settings_module:
        os.environ.setdefault(ENVIRONMENT_VARIABLE, settings_module)

    settings = LazySettings()
    apps.populate(settings.INSTALLED_APPS)

    model_data = [
        DjangoModel(
            app=app_name,
            name=model_class.__name__,
            fields=[
                ModelField(
                    name=key,
                    type_name=type(value.field).__name__,
                    opts={},
                )
                for key, value in model_class.__dict__.items()
                if isinstance(value, DeferredAttribute)
            ],
        )
        for app_name, models in apps.all_models.items()
        if models
        for model_class in models.values()
    ]

    if output:
        with output.open("w") as f:
            json.dump(model_data, f, cls=DCJSONEncoder)
        return

    json.dump(model_data, sys.stdout, cls=DCJSONEncoder)
    sys.stdout.flush()


if __name__ == "__main__":
    args = parse_args()
    run_analyzer(args.workspace_path, args.settings, args.output)
