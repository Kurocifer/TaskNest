class ToDo {
  int id;
  String todoText;
  bool isDone;

  ToDo({
    required this.id,
    required this.todoText,
    this.isDone = false,
  });

  factory ToDo.fromJson(Map<String, dynamic> json) {
    return ToDo(
      id: json['id'],
      todoText: json['task'],
      isDone: json['done'],
    );
  }

  static List<ToDo> todoList() {
    return [
      ToDo(id: 0, todoText: 'Live',),
    ];
  }
}