import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../constants/app_constants.dart';

class StorageService {
  late final SharedPreferences _prefs;

  Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
  }

  // Token management
  Future<void> saveAccessToken(String token) async {
    await _prefs.setString(AppConstants.accessTokenKey, token);
  }

  Future<String?> getAccessToken() async {
    return _prefs.getString(AppConstants.accessTokenKey);
  }

  Future<void> saveRefreshToken(String token) async {
    await _prefs.setString(AppConstants.refreshTokenKey, token);
  }

  Future<String?> getRefreshToken() async {
    return _prefs.getString(AppConstants.refreshTokenKey);
  }

  Future<void> clearTokens() async {
    await _prefs.remove(AppConstants.accessTokenKey);
    await _prefs.remove(AppConstants.refreshTokenKey);
    await _prefs.remove(AppConstants.userDataKey);
    await _prefs.setBool(AppConstants.isLoggedInKey, false);
  }

  // User data management
  Future<void> saveUserData(Map<String, dynamic> userData) async {
    await _prefs.setString(AppConstants.userDataKey, jsonEncode(userData));
  }

  Future<Map<String, dynamic>?> getUserData() async {
    final userData = _prefs.getString(AppConstants.userDataKey);
    if (userData != null) {
      return jsonDecode(userData) as Map<String, dynamic>;
    }
    return null;
  }

  // Login status
  Future<void> setLoggedIn(bool isLoggedIn) async {
    await _prefs.setBool(AppConstants.isLoggedInKey, isLoggedIn);
  }

  Future<bool> isLoggedIn() async {
    return _prefs.getBool(AppConstants.isLoggedInKey) ?? false;
  }

  // Generic storage methods
  Future<void> setString(String key, String value) async {
    await _prefs.setString(key, value);
  }

  String? getString(String key) {
    return _prefs.getString(key);
  }

  Future<void> setBool(String key, bool value) async {
    await _prefs.setBool(key, value);
  }

  bool? getBool(String key) {
    return _prefs.getBool(key);
  }

  Future<void> setInt(String key, int value) async {
    await _prefs.setInt(key, value);
  }

  int? getInt(String key) {
    return _prefs.getInt(key);
  }

  Future<void> remove(String key) async {
    await _prefs.remove(key);
  }

  Future<void> clear() async {
    await _prefs.clear();
  }
}