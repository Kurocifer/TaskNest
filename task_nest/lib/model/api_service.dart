import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'todo.dart';

class ApiService {
  final String baseUrl = 'http://localhost:8080/tasknest';

  Future<List<ToDo>> fetchTodos() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('JwtToken');

    final response = await http.get(
      Uri.parse('$baseUrl/todos/list'),
      headers: {
        'Authorization': 'Bearer $token',
      },
    );

    if (response.statusCode == 200) {
      final jsonResponse = jsonDecode(response.body);
      print('could do something');

      final List<dynamic> userTodos = jsonResponse['todos'];
      return userTodos.map((todo) => ToDo.fromJson(todo)).toList();
    } else {
      throw Exception(response.body);
    }
  }

  Future<int> deleteTodo(int todoId) async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('JwtToken');

    final url = Uri.parse('$baseUrl/todos/delete?id=$todoId');

    final response = await http.delete(
      url,
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
    );

    if (response.statusCode != 204) {
      print('Failed to delete todo: ${response.body}');
      return response.statusCode;
    }

    print('Todo deleted successfully');

    return response.statusCode;
  }

  Future<ToDo> createTodo(String task) async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('JwtToken');

    final url = Uri.parse('$baseUrl/todos/create');

    final response = await http.post(
      url,
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'task': task,
        'done': false,
      }),
    );

    if (response.statusCode != 201) {
      print('Failed to create todo: ${response.body}');
      throw Exception(
          'Failed to create todo: ${response.statusCode} - ${response.body}');
    }

    print('Todo created successfully: ${response.body}');
    final Map<String, dynamic> jsonResponse = jsonDecode(response.body);
    return ToDo.fromJson(jsonResponse);
  }

  Future<void> updateTodoStatus(int todoId, bool done) async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('JwtToken');

    final url = Uri.parse(
        '$baseUrl/todos/update-status?id=$todoId&done=${done.toString()}');

    final response = await http.put(
      url,
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
    );

    if (response.statusCode != 200) {
      print('Failed to update todo status: ${response.body}');
      throw Exception(
          'Failed to update todo status: ${response.statusCode} - ${response.body}');
    }

    print('Todo status updated successfully: ${response.body}');
  }
}
