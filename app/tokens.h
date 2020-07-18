#pragma once

#include <string>
#include <vector>
#include <memory>
#include <iostream>


using IntType = long long;

struct Token {
  virtual ~Token() = default;
};
struct IntegerToken : public Token {
  IntType value = 0;
  explicit IntegerToken(IntType value) : value(value) {}
};
struct FunctionToken : public Token {
  virtual std::unique_ptr<FunctionToken> apply(std::unique_ptr<Token>&& token) {
    std::cout << "called from base" << std::endl;
    return std::make_unique<FunctionToken>();
  };
};
struct ApplicationToken : public Token {};
struct NilToken : public Token {};
struct TokenPair : public Token {
  std::unique_ptr<Token> first;
  std::unique_ptr<Token> second;
};
struct VariableToken : public Token {
  IntType key = 0;
  explicit VariableToken(IntType key) : key(key) {};
  explicit VariableToken(std::string key) {
    if (!key.empty() && key[0] == ':') {
      key = key.substr(1);
    }
    this->key = std::stoll(key);
  }
};
struct AssignmentToken : public Token {};


void print(const Token* token) {
  if (auto int_token = dynamic_cast<const IntegerToken*>(token)) {
    std::cout << "Int: " << int_token->value << std::endl;
  }
  if (auto list_token = dynamic_cast<const ListToken*>(token)) {
    std::cout << "List token:" << std::endl;
    for (const auto& t : list_token->value) {
      print(t.get());
    }
    std::cout << std::endl;
  }
  if (auto nil_token = dynamic_cast<const NilToken*>(token)) {
    std::cout << "nil" << std::endl;
  }
}

void print(const std::unique_ptr<Token> token) {
  print(token.get());
}