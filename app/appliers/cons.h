#pragma once

#include "../tokens.h"

struct ConsFunctionToken : public FunctionToken {
  std::unique_ptr<Token> car = nullptr;
  std::unique_ptr<Token> cdr = nullptr;
  std::unique_ptr<FunctionToken> apply(std::unique_ptr<Token>&& token) override {
    auto next_p = std::make_unique<ConsFunctionToken>();
    if (!car) {
      next_p->car = std::move(token);
      return next_p;
    }
    next_p->car = std::move(car);
    next_p->cdr = std::move(token);
    return next_p;
    /*std::cout << "called right" << std::endl;
    auto next_p = std::make_unique<ConsFunctionToken>();
    std::swap(next_p->list, list);
    next_p->list.value.push_back(std::move(token));
    std::cout << "printing list" << std::endl;
    for (const auto& t : next_p->list.value) {
      print(t.get());
    }
    std::cout << "end" << std::endl;
    return next_p;*/
  }
};