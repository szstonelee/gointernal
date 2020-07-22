#include <ctime>
#include <cstdlib>
#include <vector>
#include <stack>
#include <random>
#include <algorithm>
#include <iostream>
#include <cassert>
#include <chrono>

struct Tree {
  int val;
  Tree *left;
  Tree *right;

  Tree(int v) : val(v), left(nullptr), right(nullptr) {}
};

namespace sameBST {

std::vector<int> shuffle(int n) {
  std::vector<int> nums(n);

  for (int i = 0; i < n; ++i) {
    nums[i] = i+1;
  }

  std::random_device rd;
  std::mt19937 g(rd());

  std::shuffle(nums.begin(), nums.end(), g);

  return nums;
}

Tree* insert(Tree* t, int v) {
  if (t == nullptr) return new Tree(v);

  if (v < t->val) {
    Tree* new_node = insert(t->left, v);
    if (t->left == nullptr) t->left = new_node;
    return new_node;
  } else {
    Tree* new_node = insert(t->right, v);
    if (t->right == nullptr) t->right = new_node;
    return new_node;
  }
}

Tree* buildTree(int k, int n) {
  if (k <= 0 || n <= 0) return nullptr;

  Tree* root = nullptr;
  auto nums = shuffle(n);

  for (int i = 0; i < n; ++i) {
    int v = nums[i] * k;

    Tree* new_node = insert(root, v);

    if (root == nullptr) root = new_node;
  }

  return root;
}

std::string treeToStr(const Tree* t) {
  if (t == nullptr) return "";
  
  auto l = treeToStr(t->left);
  if (!l.empty()) {
    l = "[" + l + " ";
  } else {
    l = "[";
  }

  auto r = treeToStr(t->right);
  if (!r.empty()) {
    r = " " + r + "]";
  } else {
    r = "]";
  }

  return l + std::to_string(t->val) + r;
}

void pushWholeLeft(const Tree* n, std::stack<const Tree*>* st) {
  assert(n);

  st->push(n);
  if (n->left) pushWholeLeft(n->left, st);
}

bool same(const Tree* r1, const Tree* r2) {
  std::stack<const Tree*> st1;
  std::stack<const Tree*> st2;

  if (r1) pushWholeLeft(r1, &st1);
  if (r2) pushWholeLeft(r2, &st2);

  while (!st1.empty() && !st2.empty()) {
    auto n1 = st1.top();
    st1.pop();
    auto n2 = st2.top();
    st2.pop();

    if (n1->val != n2->val) return false;

    if (n1->right) pushWholeLeft(n1->right, &st1);
    if (n2->right) pushWholeLeft(n2->right, &st2);
  }
  return st1.empty() && st2.empty();
}

} // end namespace

int main(int argc, char** argv) {

  auto start1 = std::chrono::steady_clock::now();
  Tree* r1 = sameBST::buildTree(1, 1 << 23);
  Tree* r2 = sameBST::buildTree(1, 1 << 23);
  auto end1 = std::chrono::steady_clock::now();
  std::cout << "Tree build time = " << 
    std::chrono::duration_cast<std::chrono::seconds>(end1-start1).count() << "(s)" << std::endl;

  auto start2 = std::chrono::steady_clock::now();
  auto is_same = sameBST::same(r1, r2);
  auto end2 = std::chrono::steady_clock::now();
  std::cout << "BST Tree " << (is_same ? "same" : "not same") << std::endl;
  std::cout << "Elaspsed time for compute same = " << 
    std::chrono::duration_cast<std::chrono::milliseconds>(end2-start2).count() << "(ms)" << std::endl;

  return 0;
}

