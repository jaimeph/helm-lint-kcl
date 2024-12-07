schema Values:
    podspec: PodSpec

schema PodSpec:
    containers: {str:Container}

schema Container:
    enabled: bool = True
    image: Image

schema Image:
    repository: str
    pullPolicy?: "IfNotPresent" | "Always" | "Never" = "IfNotPresent"
    tag: str
