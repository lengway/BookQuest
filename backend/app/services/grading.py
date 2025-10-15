from typing import List, Dict, Tuple
from ..models import Question, Option, QuestionType


def grade_single_multi(selected_option_ids: List[int], options: List[Option], multi: bool) -> bool:
    correct_ids = {opt.id for opt in options if opt.is_correct}
    selected_ids = set(selected_option_ids or [])
    if not multi and len(selected_ids) != 1:
        return False
    return selected_ids == correct_ids


def grade_ordering(ordering: List[int], options: List[Option]) -> bool:
    expected = [opt.id for opt in sorted(options, key=lambda o: (o.order_index or 0))]
    return ordering == expected


def grade_matching(matches: Dict[str, str], options: List[Option]) -> bool:
    # Для простоты: считаем, что у опций есть match_key, пары верны если ключи совпали по ожиданию.
    # Ожидаем соответствие: левый элемент (у которого match_key не None) → правый с тем же match_key.
    # Матч выражен как left_key -> right_key.
    if not matches:
        return False
    # соберём множество пар из опций
    keys = [opt.match_key for opt in options if opt.match_key is not None]
    unique_keys = set(keys)
    # правильные пары: k -> k для каждого k
    for k, v in matches.items():
        if k not in unique_keys or v not in unique_keys or k != v:
            return False
    return True


def grade_question(question: Question, submission_payload: dict) -> bool:
    qtype = question.type
    options = question.options
    if qtype == QuestionType.SINGLE_CHOICE:
        return grade_single_multi(submission_payload.get("selected_option_ids", []), options, multi=False)
    if qtype == QuestionType.MULTI_CHOICE:
        return grade_single_multi(submission_payload.get("selected_option_ids", []), options, multi=True)
    if qtype == QuestionType.ORDERING:
        return grade_ordering(submission_payload.get("ordering", []), options)
    if qtype == QuestionType.MATCHING:
        return grade_matching(submission_payload.get("matches", {}), options)
    return False


