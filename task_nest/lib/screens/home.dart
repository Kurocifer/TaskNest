import 'package:flutter/material.dart';
import 'package:task_nest/components/color_button.dart';
import 'package:task_nest/components/theme_button.dart';
import 'package:task_nest/constants/constants.dart';
import 'package:task_nest/model/api_service.dart';
import 'package:task_nest/model/auth.dart';
import 'package:task_nest/screens/account_details_page.dart';
import 'package:go_router/go_router.dart';

import '../model/todo.dart';
import '../widgets/todo_item.dart';

String _notFoundMessage =
    ""; // the message to be returned when search matches no existing todo
Size? screenSize;

class Home extends StatefulWidget {
  final void Function(bool useLightMode) changeTheme;
  final void Function(int value) changeColor;
  final ColorSelection colorSelected;
  final int tab;
  final TaskNestAuth auth;

  const Home({
    super.key,
    required this.changeTheme,
    required this.changeColor,
    required this.colorSelected,
    required this.tab,
    required this.auth,
  });

  @override
  State<Home> createState() => _HomeState();
}

class _HomeState extends State<Home> {
  String filter = 'all'; // Show all Todos by default
  List<ToDo> _foundToDo = [];
  List<ToDo> todosList = [];
  final _todoController = TextEditingController();
  final GlobalKey<ScaffoldState> scaffoldKey = GlobalKey<ScaffoldState>();
  static const double drawerWidth = 375.0;
  List<ToDo> filteredTodos = [];
  final FocusNode _focusNode = FocusNode();
  final ApiService apiService = ApiService();

  @override
  void initState() {
    _foundToDo = todosList;
    super.initState();
    _fetchTodos();
  }

  @override
  Widget build(BuildContext context) {
    // Filter todos based on the value of filter

    return Scaffold(
      backgroundColor: Theme.of(context).colorScheme.surface,
      key: scaffoldKey,
      drawer: _buildDrawer(),
      appBar: _buildAppBar(),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    final textTheme = Theme.of(context)
        .textTheme
        .apply(displayColor: Theme.of(context).colorScheme.onSurface);

    return Stack(
      children: [
        Container(
          padding: const EdgeInsets.symmetric(
            horizontal: 20,
            vertical: 15,
          ),
          child: Column(
            children: [
              searchBox(),
              Expanded(
                child: ListView(
                  children: [
                    Container(
                      margin: const EdgeInsets.only(
                        top: 50,
                        bottom: 20,
                      ),
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      children: [
                        TextButton(
                          onPressed: () {
                            setState(() {
                              filter = 'all';
                            });
                            _runTodoFiltering();
                          },
                          child: Text(
                            'All Todos',
                            style: TextStyle(
                                fontSize: 30,
                                fontWeight: FontWeight.w500,
                                decoration: filter == 'all'
                                    ? TextDecoration.underline
                                    : TextDecoration.none),
                          ),
                        ),
                        TextButton(
                          onPressed: () {
                            setState(() {
                              filter = 'done';
                            });
                            _runTodoFiltering();
                          },
                          child: Text(
                            'Done Todos',
                            style: TextStyle(
                              fontSize: 30.0,
                              fontWeight: FontWeight.w500,
                              decoration: filter == 'done'
                                  ? TextDecoration.underline
                                  : TextDecoration.none,
                            ),
                          ),
                        ),
                        TextButton(
                          onPressed: () {
                            setState(() {
                              filter = 'undone';
                            });
                            _runTodoFiltering();
                          },
                          child: Text(
                            'Undone Todos',
                            style: TextStyle(
                              fontSize: 30,
                              fontWeight: FontWeight.w500,
                              decoration: filter == 'undone'
                                  ? TextDecoration.underline
                                  : TextDecoration.none,
                            ),
                          ),
                        ),
                      ],
                    ),
                    Container(
                      margin: const EdgeInsets.only(
                        top: 10,
                        bottom: 20,
                      ),
                    ),
                    if (_notFoundMessage.isNotEmpty)
                      Center(
                        child: Text(
                          _notFoundMessage,
                          style: textTheme.titleLarge,
                        ),
                      ),
                    for (ToDo todoItem in _foundToDo.reversed)
                      ToDoItem(
                        todo: todoItem,
                        onToDoChanged: _handleToDoChange,
                        onDeleteItem: _deleteToDoItem,
                      ),
                  ],
                ),
              )
            ],
          ),
        ),
        Align(
          alignment: Alignment.bottomCenter,
          child: Row(children: [
            Expanded(
              child: Card(
                  elevation: 10.0,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10.0),
                  ),
                  child: TextField(
                    autofocus:
                        true, // get focus on the textfield when widget is initialized
                    focusNode: _focusNode,
                    controller: _todoController,
                    onSubmitted: _addToDoItem,
                    decoration: const InputDecoration(
                        hintText: 'Add a new todo item',
                        border: InputBorder.none,
                        contentPadding: EdgeInsets.all(16.0)),
                  )),
            ),
            Container(
              margin: const EdgeInsets.only(
                bottom: 20,
                right: 20,
              ),
              child: ElevatedButton(
                onPressed: () {
                  _addToDoItem(_todoController.text);
                },
                style: ElevatedButton.styleFrom(
                  minimumSize: const Size(60, 60),
                  elevation: 10.0,
                ),
                child: const Text(
                  '+',
                  style: TextStyle(
                    fontSize: 40,
                  ),
                ),
              ),
            ),
          ]),
        ),
      ],
    );
  }

  Future<void> _fetchTodos() async {
    try {
      final List<ToDo> fetchedTodos = await apiService.fetchTodos();
      setState(() {
        todosList = fetchedTodos;
        _foundToDo = fetchedTodos;
      });
    } catch (e) {
      print('Error fetching todos: $e');
    }
  }

  void openDrawer() {
    scaffoldKey.currentState!.openDrawer();
  }

  void _handleToDoChange(ToDo todo) async {
    try {
      await apiService.updateTodoStatus(
          todo.id, !todo.isDone); // Toggle the status
      print('Todo status updated for ID: ${todo.id}');
      setState(() {
        todo.isDone = !todo.isDone;
        _runTodoFiltering(); // call this here so when viewing done todos and a todo is unmarked, it should not longer be viewed.
      });
    } catch (e) {
      print('Error updating todo status: $e');
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to update todo status: $e')),
      );
    }
  }

  void _deleteToDoItem(int id) async {
    int status = await apiService.deleteTodo(id);
    setState(() {
      if (status == 204) {
        todosList.removeWhere((item) => item.id == id);
      }
    });
  }

  void _addToDoItem(String todoTask) async {
    if (todoTask.isNotEmpty) {
      try {
        final createdTodo = await apiService.createTodo(todoTask);
        print(
            'Todo created: ${createdTodo.todoText} with ID: ${createdTodo.id}');
        setState(() {
          todosList.add(createdTodo);
        });
      } catch (e) {
        print('Error creating todo: $e');
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to create todo: $e')),
        );
      }
    } else {
      print('Task cannot be empty');
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Task cannot be empty')),
      );
    }

    _todoController.clear();
    _focusNode.requestFocus(); // return focus to the text field.
  }

  void _runFilter(String enteredKeyword) {
    List<ToDo> results = [];
    String message = "";

    if (enteredKeyword.isEmpty) {
      results = todosList;
    } else {
      results = filteredTodos
          .where((item) => item.todoText
              .toLowerCase()
              .contains(enteredKeyword.toLowerCase()))
          .toList();
    }

    // Check if results are empty and set the message
    if (results.isEmpty) {
      message = "Todo not found";
    }

    setState(() {
      _foundToDo = results;
      _notFoundMessage = message;
    });
  }

  void _runTodoFiltering() {
    _notFoundMessage = '';
    filteredTodos = todosList.where((todo) {
      if (filter == 'done') return todo.isDone; // return only done todos
      if (filter == 'undone') return !todo.isDone; // return only undone todos
      return true; // return all todos
    }).toList();

    if (filteredTodos.isEmpty) {
      switch (filter) {
        case 'done':
          _notFoundMessage = 'No done todos yet';
          break;
        case 'undone':
          _notFoundMessage = 'No undone todos';
          break;

        case 'all':
          _notFoundMessage = 'No todos yet';
      }
    }
    setState(() {
      _foundToDo = filteredTodos;
    });
  }

  Widget searchBox() {
    return Card(
      elevation: 2.0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(10.0),
      ),
      child: TextField(
        onChanged: (value) => _runFilter(value),
        decoration: const InputDecoration(
          hintText: 'Search...',
          border: InputBorder.none,
          contentPadding: EdgeInsets.all(16.0),
          prefixIcon: Icon(Icons.search),
        ),
      ),
    );
  }

  Widget _buildDrawer() {
    return SizedBox(
      width: drawerWidth,
      child: Drawer(
        child: AccountDetailsPage(
          onLogOut: (logout) async {
            widget.auth.signOut().then((value) => context.go('/login'));
          },
        ),
      ),
    );
  }

  AppBar _buildAppBar() {
    return AppBar(
      backgroundColor: Theme.of(context).colorScheme.surface,
      elevation: 4.0,
      leading: IconButton(
        icon: SizedBox(
          height: 40.0,
          width: 40.0,
          child: ClipRRect(
              borderRadius: BorderRadius.circular(20.0),
              child: Image.asset('assets/images/cifer.png')),
        ),
        onPressed: openDrawer,
      ),
      title: Row(
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          Row(
            children: [
              ColorButton(
                changeColor: widget.changeColor,
                colorSelected: widget.colorSelected,
              ),
              ThemeButton(
                changeThemeMode: widget.changeTheme,
              ),
            ],
          )
        ],
      ),
    );
  }
}
