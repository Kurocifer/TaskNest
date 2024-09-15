import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:task_nest/model/app_cache.dart';

class AccountDetailsPage extends StatefulWidget {
  final Function onLogOut;

  const AccountDetailsPage({
    super.key,
    required this.onLogOut,
  });

  @override
  State<AccountDetailsPage> createState() => _AccountDetailsPageState();
}

class _AccountDetailsPageState extends State<AccountDetailsPage> {
  final appCache = AppCache();
  @override
  Widget build(BuildContext context) {

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: FutureBuilder(
        future: _getUserData(),
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(
              child: CircularProgressIndicator(),
            );
          }

          if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          }

          final userData = snapshot.data as Map<String, String?>;

          return Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                SizedBox(
                  height: 70.0,
                  width: 70.0,
                  child: ClipRRect(
                      borderRadius: BorderRadius.circular(40.0),
                      child: Image.asset('assets/images/cifer.png')),
                ),
                Text(
                  '${userData['username']}',
                    style: const TextStyle(fontSize: 20)),
                const SizedBox(height: 10),
                Center(
              child: ListTile(
                title: const Text(
                  'Log out',
                  style: TextStyle(fontSize: 30),
                ),
                onTap: () {
                  widget.onLogOut(true);
                },
              ),
            ),
              ],
            ),
          );
        },
      ),
    );
  }
}

Future<Map<String, String?>> _getUserData() async {
  final prefs = await SharedPreferences.getInstance();
  final token = prefs.getString('jwtToken');
  final userId = prefs.getInt('userId');
  final username = prefs.getString('username');

  return {
    'token': token,
    'userId': userId.toString(),
    'username': username,
  };
}
