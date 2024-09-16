import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:task_nest/model/todo.dart';

class FetchTodos with ChangeNotifier {
  List<ToDo> _todos = [];
  bool _isLoading = false;

  List<ToDo> get todos => _todos;
  bool get isLoading => _isLoading;
  Future<void> fetchTodos(String userId) async {
    _isLoading = true; // Start loading
    notifyListeners(); // Notify listeners about the loading state

  final prefs = await SharedPreferences.getInstance();
  final token = prefs.getString('jwtToken');; // Retrieve the JWT token

    const url = 'https://localhost:8080/tasknest/todos/list';

    final response = await http.get(
      Uri.parse(url),
      headers: {
        'Authorization': 'Bearer $token', // Include the token in the headers
      },
    );

    if (response.statusCode == 200) {
      final List<dynamic> todoList = jsonDecode(response.body);
      _todos = todoList.map((json) => ToDo.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load todos');
    }

    _isLoading = false; // Stop loading
    notifyListeners(); // Notify listeners again after loading
  }
}