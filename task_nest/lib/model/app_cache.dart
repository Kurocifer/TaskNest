import 'package:shared_preferences/shared_preferences.dart';

class AppCache {
  static const klogin = 'task_nest_login';
  static const ktoken  = 'JwtToken';
  static const kUsername = 'username';
  static const kUserId = 'userId';

  Future<void> invalidate() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool(klogin, false);
    await prefs.remove(ktoken);
    await prefs.remove(kUserId);
    await prefs.remove(kUsername);
  }

  Future<void> cacheUser(String token, String username, int userId) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool(klogin, true);
    await prefs.setString(ktoken, token);
    await prefs.setInt(kUserId, userId);
    await prefs.setString(kUsername, username);
  }

  Future<bool> isUserLoggedIn() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getBool(klogin) ?? false;
  }
}