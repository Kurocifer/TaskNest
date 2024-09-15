import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;

import 'app_cache.dart';

/// A mock authentication service.
class TaskNestAuth extends ChangeNotifier {
  bool _loggedIn = false;

  // Stores user state properties on platform specific file system.
  final _appCache = AppCache();

  Future<bool> get loggedIn => _appCache.isUserLoggedIn();

  /// Signs out the current user.
  Future<void> signOut() async {
    await Future<void>.delayed(const Duration(milliseconds: 200));
    // Sign out.
    _loggedIn = false;
    await _appCache.invalidate();
    notifyListeners();
  }

  /// Signs in a user.
  Future<bool> signIn(
      BuildContext context, String username, String password) async {
    const url = 'http://localhost:8080/tasknest/login';

    final response = await http.post(
      Uri.parse(url),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: jsonEncode(<String, String>{
        "username": username,
        "password": password,
      }),
    );

    final responseData = jsonDecode(response.body);

    if (response.statusCode == 200) {
      final token = responseData['token'];
      final username = responseData['username'];
      final userId = responseData['id'];

      await _appCache.cacheUser(token, username, userId);
      _loggedIn = true;
      notifyListeners();
      return _loggedIn;
    } else {
      String errorMessage = responseData['error'];
      _loggedIn = false;
      _appCache.invalidate();
      notifyListeners();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('An error occured: $errorMessage')),
      );
      throw Exception('Failed to login: ${response.body}');
    }
  }

  Future<bool> signUp(
      BuildContext context, String username, String password) async {
    const url = 'http://localhost:8080/tasknest/register';

    final response = await http.post(
      Uri.parse(url),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: jsonEncode(<String, String>{
        "username": username,
        "password": password,
      }),
    );

    final responseData = jsonDecode(response.body);

    if (response.statusCode == 201) {
      final token = responseData['token'];
      final username = responseData['username'];
      final userId = responseData['id'];

      await _appCache.cacheUser(token, username, userId);
      _loggedIn = true;
      notifyListeners();
      return _loggedIn;
    } else {
      String errorMessage = responseData['error'];
      _loggedIn = false;
      _appCache.invalidate();
      notifyListeners();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('An error occured: $errorMessage')),
      );
      throw Exception('Failed to login: ${response.body}');
    }
  }
}
