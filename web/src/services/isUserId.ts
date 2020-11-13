export function isUserId(value: string) {
    return value.startsWith("id:");
}

export function getUserId(value: string) {
    return value.replace("id:", "");
}