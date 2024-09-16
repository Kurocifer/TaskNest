import 'package:json_annotation/json_annotation.dart';

//part 'user_model.g.dart';

@JsonSerializable()
class User {
  final String id;
  final String username;
  String? profileImageUrl;
  String? preferredColor;
  bool? darkmode;

  User({
    required this.id,
    required this.username,
    this.profileImageUrl,
    this.preferredColor,
    this.darkmode,
  }) {
    if(profileImageUrl == null || profileImageUrl!.isEmpty) {
      profileImageUrl = 'assets/images/cifer.png';
    }

    if(preferredColor == null || preferredColor!.isEmpty) {
      preferredColor = 'White';
    }

    darkmode ??= false;
  }

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'],
      username: json['username']
    );
  }

}
