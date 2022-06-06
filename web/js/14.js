(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([[14],{"./node_modules/cache-loader/dist/cjs.js?!./node_modules/babel-loader/lib/index.js!./node_modules/cache-loader/dist/cjs.js?!./node_modules/vue-loader/lib/index.js?!./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=script&lang=js&":
/*!****************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/cache-loader/dist/cjs.js??ref--12-0!./node_modules/babel-loader/lib!./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options!./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=script&lang=js& ***!
  \****************************************************************************************************************************************************************************************************************************************************************************/
/*! exports provided: default */function(module,__webpack_exports__,__webpack_require__){"use strict";eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var core_js_modules_es_array_filter_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! core-js/modules/es.array.filter.js */ \"./node_modules/core-js/modules/es.array.filter.js\");\n/* harmony import */ var core_js_modules_es_array_filter_js__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(core_js_modules_es_array_filter_js__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var core_js_modules_web_dom_collections_for_each_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! core-js/modules/web.dom-collections.for-each.js */ \"./node_modules/core-js/modules/web.dom-collections.for-each.js\");\n/* harmony import */ var core_js_modules_web_dom_collections_for_each_js__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(core_js_modules_web_dom_collections_for_each_js__WEBPACK_IMPORTED_MODULE_1__);\n/* harmony import */ var _mixins_mixin__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @/mixins/mixin */ \"./src/mixins/mixin.js\");\n/* harmony import */ var viewerjs_dist_viewer_css__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! viewerjs/dist/viewer.css */ \"./node_modules/viewerjs/dist/viewer.css\");\n/* harmony import */ var viewerjs_dist_viewer_css__WEBPACK_IMPORTED_MODULE_3___default = /*#__PURE__*/__webpack_require__.n(viewerjs_dist_viewer_css__WEBPACK_IMPORTED_MODULE_3__);\n/* harmony import */ var v_viewer__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! v-viewer */ \"./node_modules/v-viewer/dist/v-viewer.js\");\n/* harmony import */ var v_viewer__WEBPACK_IMPORTED_MODULE_4___default = /*#__PURE__*/__webpack_require__.n(v_viewer__WEBPACK_IMPORTED_MODULE_4__);\n\n\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n//\n\n\n\nvar XIMAGES = ['png', 'jpg', 'jpeg', 'bmp', 'gif', 'webp', 'svg', 'tiff'];\n/* harmony default export */ __webpack_exports__[\"default\"] = ({\n  mixins: [_mixins_mixin__WEBPACK_IMPORTED_MODULE_2__[\"mixin\"]],\n  props: {\n    item: {\n      type: Object,\n      default: function _default() {\n        return {\n          path: '',\n          name: ''\n        };\n      }\n    },\n    list: []\n  },\n  components: {\n    Viewer: v_viewer__WEBPACK_IMPORTED_MODULE_4__[\"component\"]\n  },\n  data: function data() {\n    return {\n      isMoving: false,\n      timeout: null,\n      itemList: [],\n      currentItem: this.item,\n      currentItemIndex: 0,\n      currentItemArray: [],\n      viewer: {},\n      viewerOptions: {\n        button: false,\n        //Hide FullScreen Button\n        toolbar: false,\n        //Hide Toolbar\n        title: false,\n        //Hide Title\n        navbar: false,\n        //Hide Navbar\n        backdrop: false,\n        //Hide Background\n        transition: false,\n        //Without css3 animation\n        inline: true,\n        initialViewIndex: 0\n      }\n    };\n  },\n  computed: {\n    disableNext: function disableNext() {\n      return this.currentItemIndex == this.itemList.length - 1;\n    },\n    disablePrev: function disablePrev() {\n      return this.currentItemIndex == 0;\n    }\n  },\n  created: function created() {\n    this.filterImages();\n    this.getCurrentImageIndex();\n    this.setSourceImageURLs();\n  },\n  mounted: function mounted() {\n    var _this = this;\n\n    window.onkeyup = function (e) {\n      switch (e.code) {\n        case 'ArrowRight':\n          _this.next();\n\n          break;\n\n        case 'ArrowLeft':\n          _this.prev();\n\n          break;\n      }\n    };\n  },\n  methods: {\n    download: function download() {\n      this.downloadFile(this.currentItem);\n    },\n    close: function close() {\n      this.$emit(\"close\");\n    },\n    inited: function inited(viewer) {\n      this.viewer = viewer;\n      this.viewer.show();\n      this.onMouseMove();\n    },\n    next: function next() {\n      if (this.currentItemIndex < this.itemList.length - 1) {\n        this.currentItemIndex++;\n        this.setSourceImageURLs();\n      }\n    },\n    prev: function prev() {\n      if (this.currentItemIndex > 0) {\n        this.currentItemIndex--;\n        this.setSourceImageURLs();\n      }\n    },\n    filterImages: function filterImages() {\n      var _this2 = this;\n\n      this.itemList = this.list.filter(function (item) {\n        var ext = _this2.getFileExt(item);\n\n        return !item.is_dir && XIMAGES.indexOf(ext.toLowerCase()) > -1;\n      });\n    },\n    getCurrentImageIndex: function getCurrentImageIndex() {\n      var _this3 = this;\n\n      this.itemList.forEach(function (item, index) {\n        if (item == _this3.currentItem) {\n          _this3.currentItemIndex = index;\n        }\n      });\n    },\n    setSourceImageURLs: function setSourceImageURLs() {\n      this.currentItem = this.itemList[this.currentItemIndex];\n      this.currentItemArray = [this.getFileUrl(this.currentItem)];\n    },\n    // Hide Toolbar after 5 seconds\n    onMouseMove: function onMouseMove() {\n      var _this4 = this;\n\n      this.isMoving = true;\n\n      if (this.timeout !== null) {\n        clearTimeout(this.timeout);\n      }\n\n      this.timeout = setTimeout(function () {\n        _this4.isMoving = false;\n        _this4.timeout = null;\n      }, 5000);\n    }\n  }\n});\n\n//# sourceURL=webpack:///./src/components/filebrowser/viewers/ImageViewer.vue?./node_modules/cache-loader/dist/cjs.js??ref--12-0!./node_modules/babel-loader/lib!./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options")},'./node_modules/cache-loader/dist/cjs.js?{"cacheDirectory":"node_modules/.cache/vue-loader","cacheIdentifier":"d12f3824-vue-loader-template"}!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/cache-loader/dist/cjs.js?!./node_modules/vue-loader/lib/index.js?!./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=template&id=5973f308&':
/*!************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** ./node_modules/cache-loader/dist/cjs.js?{"cacheDirectory":"node_modules/.cache/vue-loader","cacheIdentifier":"d12f3824-vue-loader-template"}!./node_modules/vue-loader/lib/loaders/templateLoader.js??vue-loader-options!./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options!./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=template&id=5973f308& ***!
  \************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************/
/*! exports provided: render, staticRenderFns */function(module,__webpack_exports__,__webpack_require__){"use strict";eval('__webpack_require__.r(__webpack_exports__);\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "render", function() { return render; });\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return staticRenderFns; });\nvar render = function() {\n  var _vm = this\n  var _h = _vm.$createElement\n  var _c = _vm._self._c || _h\n  return _c(\n    "div",\n    {\n      staticClass: "overlay",\n      attrs: { id: "image_viewer" },\n      on: { mousemove: _vm.onMouseMove, touchmove: _vm.onMouseMove }\n    },\n    [\n      _c("header", { staticClass: "modal-card-head" }, [\n        _c("div", { staticClass: "is-flex  is-flex-grow-1 is-flex-shrink-1" }, [\n          _c("h3", { staticClass: "title is-6 one-line" }, [\n            _vm._v(_vm._s(_vm.currentItem.name))\n          ])\n        ]),\n        _c(\n          "div",\n          { staticClass: "is-flex is-align-items-center is-flex-shrink-0" },\n          [\n            _c("b-button", {\n              staticClass: "mr-2",\n              attrs: {\n                "icon-left": "download",\n                type: "is-primary",\n                size: "is-small",\n                label: _vm.$t("Download"),\n                rounded: ""\n              },\n              on: { click: _vm.download }\n            }),\n            _c(\n              "div",\n              {\n                staticClass:\n                  "is-flex is-align-items-center modal-close-container close-btn modal-close-container-line"\n              },\n              [\n                _c("button", {\n                  staticClass: "delete",\n                  attrs: { type: "button" },\n                  on: { click: _vm.close }\n                })\n              ]\n            )\n          ],\n          1\n        )\n      ]),\n      _vm.isMoving\n        ? _c(\n            "div",\n            { staticClass: "img-toolbar is-flex" },\n            [\n              _c(\n                "b-tooltip",\n                { attrs: { label: _vm.$t("Previous"), type: "is-dark" } },\n                [\n                  _c(\n                    "div",\n                    {\n                      staticClass: "toolbar-item",\n                      class: { disabled: _vm.disablePrev },\n                      on: { click: _vm.prev }\n                    },\n                    [\n                      _c(\n                        "span",\n                        { staticClass: "has-text-white block" },\n                        [\n                          _c("b-icon", {\n                            attrs: { icon: "arrow-left-thin", size: "is-small" }\n                          })\n                        ],\n                        1\n                      )\n                    ]\n                  )\n                ]\n              ),\n              _c(\n                "b-tooltip",\n                { attrs: { label: _vm.$t("Zoom in"), type: "is-dark" } },\n                [\n                  _c(\n                    "div",\n                    {\n                      staticClass: "toolbar-item",\n                      on: {\n                        click: function($event) {\n                          return _vm.viewer.zoom(0.1)\n                        }\n                      }\n                    },\n                    [\n                      _c(\n                        "span",\n                        { staticClass: "has-text-white block" },\n                        [\n                          _c("b-icon", {\n                            attrs: {\n                              icon: "magnify-plus-outline",\n                              size: "is-small"\n                            }\n                          })\n                        ],\n                        1\n                      )\n                    ]\n                  )\n                ]\n              ),\n              _c(\n                "b-tooltip",\n                { attrs: { label: _vm.$t("Rotate"), type: "is-dark" } },\n                [\n                  _c(\n                    "div",\n                    {\n                      staticClass: "toolbar-item",\n                      on: {\n                        click: function($event) {\n                          return _vm.viewer.rotate(90)\n                        }\n                      }\n                    },\n                    [\n                      _c(\n                        "span",\n                        { staticClass: "has-text-white block" },\n                        [\n                          _c("b-icon", {\n                            attrs: {\n                              icon: "format-rotate-90",\n                              size: "is-small",\n                              "custom-class": "mdi-flip-h"\n                            }\n                          })\n                        ],\n                        1\n                      )\n                    ]\n                  )\n                ]\n              ),\n              _c(\n                "b-tooltip",\n                { attrs: { label: _vm.$t("Reset"), type: "is-dark" } },\n                [\n                  _c(\n                    "div",\n                    {\n                      staticClass: "toolbar-item",\n                      on: {\n                        click: function($event) {\n                          return _vm.viewer.reset()\n                        }\n                      }\n                    },\n                    [\n                      _c(\n                        "span",\n                        { staticClass: "has-text-white block" },\n                        [\n                          _c("b-icon", {\n                            attrs: { icon: "restore", size: "is-small" }\n                          })\n                        ],\n                        1\n                      )\n                    ]\n                  )\n                ]\n              ),\n              _c(\n                "b-tooltip",\n                { attrs: { label: _vm.$t("Zoom out"), type: "is-dark" } },\n                [\n                  _c(\n                    "div",\n                    {\n                      staticClass: "toolbar-item",\n                      on: {\n                        click: function($event) {\n                          return _vm.viewer.zoom(-0.1)\n                        }\n                      }\n                    },\n                    [\n                      _c(\n                        "span",\n                        { staticClass: "has-text-white block" },\n                        [\n                          _c("b-icon", {\n                            attrs: {\n                              icon: "magnify-minus-outline",\n                              size: "is-small"\n                            }\n                          })\n                        ],\n                        1\n                      )\n                    ]\n                  )\n                ]\n              ),\n              _c(\n                "b-tooltip",\n                { attrs: { label: _vm.$t("INext"), type: "is-dark" } },\n                [\n                  _c(\n                    "div",\n                    {\n                      staticClass: "toolbar-item",\n                      class: { disabled: _vm.disableNext },\n                      on: { click: _vm.next }\n                    },\n                    [\n                      _c(\n                        "span",\n                        { staticClass: "has-text-white block" },\n                        [\n                          _c("b-icon", {\n                            attrs: {\n                              icon: "arrow-right-thin",\n                              size: "is-small"\n                            }\n                          })\n                        ],\n                        1\n                      )\n                    ]\n                  )\n                ]\n              )\n            ],\n            1\n          )\n        : _vm._e(),\n      _c(\n        "div",\n        { staticClass: " v-container pl-4 pr-4" },\n        [\n          _c("viewer", {\n            ref: "viewer",\n            staticClass: "viewer",\n            attrs: { options: _vm.viewerOptions, images: _vm.currentItemArray },\n            on: { inited: _vm.inited },\n            scopedSlots: _vm._u([\n              {\n                key: "default",\n                fn: function(scope) {\n                  return _vm._l(scope.images, function(src) {\n                    return _c("img", { key: src, attrs: { src: src } })\n                  })\n                }\n              }\n            ])\n          })\n        ],\n        1\n      )\n    ]\n  )\n}\nvar staticRenderFns = []\nrender._withStripped = true\n\n\n\n//# sourceURL=webpack:///./src/components/filebrowser/viewers/ImageViewer.vue?./node_modules/cache-loader/dist/cjs.js?%7B%22cacheDirectory%22:%22node_modules/.cache/vue-loader%22,%22cacheIdentifier%22:%22d12f3824-vue-loader-template%22%7D!./node_modules/vue-loader/lib/loaders/templateLoader.js??vue-loader-options!./node_modules/cache-loader/dist/cjs.js??ref--0-0!./node_modules/vue-loader/lib??vue-loader-options')},"./src/components/filebrowser/viewers/ImageViewer.vue":
/*!************************************************************!*\
  !*** ./src/components/filebrowser/viewers/ImageViewer.vue ***!
  \************************************************************/
/*! exports provided: default */function(module,__webpack_exports__,__webpack_require__){"use strict";eval('__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _ImageViewer_vue_vue_type_template_id_5973f308___WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./ImageViewer.vue?vue&type=template&id=5973f308& */ "./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=template&id=5973f308&");\n/* harmony import */ var _ImageViewer_vue_vue_type_script_lang_js___WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./ImageViewer.vue?vue&type=script&lang=js& */ "./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=script&lang=js&");\n/* empty/unused harmony star reexport *//* harmony import */ var _node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../../../../node_modules/vue-loader/lib/runtime/componentNormalizer.js */ "./node_modules/vue-loader/lib/runtime/componentNormalizer.js");\n\n\n\n\n\n/* normalize component */\n\nvar component = Object(_node_modules_vue_loader_lib_runtime_componentNormalizer_js__WEBPACK_IMPORTED_MODULE_2__["default"])(\n  _ImageViewer_vue_vue_type_script_lang_js___WEBPACK_IMPORTED_MODULE_1__["default"],\n  _ImageViewer_vue_vue_type_template_id_5973f308___WEBPACK_IMPORTED_MODULE_0__["render"],\n  _ImageViewer_vue_vue_type_template_id_5973f308___WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"],\n  false,\n  null,\n  null,\n  null\n  \n)\n\n/* hot reload */\nif (false) { var api; }\ncomponent.options.__file = "src/components/filebrowser/viewers/ImageViewer.vue"\n/* harmony default export */ __webpack_exports__["default"] = (component.exports);\n\n//# sourceURL=webpack:///./src/components/filebrowser/viewers/ImageViewer.vue?')},"./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=script&lang=js&":
/*!*************************************************************************************!*\
  !*** ./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=script&lang=js& ***!
  \*************************************************************************************/
/*! exports provided: default */function(module,__webpack_exports__,__webpack_require__){"use strict";eval('__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _node_modules_cache_loader_dist_cjs_js_ref_12_0_node_modules_babel_loader_lib_index_js_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_ImageViewer_vue_vue_type_script_lang_js___WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/cache-loader/dist/cjs.js??ref--12-0!../../../../node_modules/babel-loader/lib!../../../../node_modules/cache-loader/dist/cjs.js??ref--0-0!../../../../node_modules/vue-loader/lib??vue-loader-options!./ImageViewer.vue?vue&type=script&lang=js& */ "./node_modules/cache-loader/dist/cjs.js?!./node_modules/babel-loader/lib/index.js!./node_modules/cache-loader/dist/cjs.js?!./node_modules/vue-loader/lib/index.js?!./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=script&lang=js&");\n/* empty/unused harmony star reexport */ /* harmony default export */ __webpack_exports__["default"] = (_node_modules_cache_loader_dist_cjs_js_ref_12_0_node_modules_babel_loader_lib_index_js_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_ImageViewer_vue_vue_type_script_lang_js___WEBPACK_IMPORTED_MODULE_0__["default"]); \n\n//# sourceURL=webpack:///./src/components/filebrowser/viewers/ImageViewer.vue?')},"./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=template&id=5973f308&":
/*!*******************************************************************************************!*\
  !*** ./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=template&id=5973f308& ***!
  \*******************************************************************************************/
/*! exports provided: render, staticRenderFns */function(module,__webpack_exports__,__webpack_require__){"use strict";eval('__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _node_modules_cache_loader_dist_cjs_js_cacheDirectory_node_modules_cache_vue_loader_cacheIdentifier_d12f3824_vue_loader_template_node_modules_vue_loader_lib_loaders_templateLoader_js_vue_loader_options_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_ImageViewer_vue_vue_type_template_id_5973f308___WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! -!../../../../node_modules/cache-loader/dist/cjs.js?{"cacheDirectory":"node_modules/.cache/vue-loader","cacheIdentifier":"d12f3824-vue-loader-template"}!../../../../node_modules/vue-loader/lib/loaders/templateLoader.js??vue-loader-options!../../../../node_modules/cache-loader/dist/cjs.js??ref--0-0!../../../../node_modules/vue-loader/lib??vue-loader-options!./ImageViewer.vue?vue&type=template&id=5973f308& */ "./node_modules/cache-loader/dist/cjs.js?{\\"cacheDirectory\\":\\"node_modules/.cache/vue-loader\\",\\"cacheIdentifier\\":\\"d12f3824-vue-loader-template\\"}!./node_modules/vue-loader/lib/loaders/templateLoader.js?!./node_modules/cache-loader/dist/cjs.js?!./node_modules/vue-loader/lib/index.js?!./src/components/filebrowser/viewers/ImageViewer.vue?vue&type=template&id=5973f308&");\n/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "render", function() { return _node_modules_cache_loader_dist_cjs_js_cacheDirectory_node_modules_cache_vue_loader_cacheIdentifier_d12f3824_vue_loader_template_node_modules_vue_loader_lib_loaders_templateLoader_js_vue_loader_options_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_ImageViewer_vue_vue_type_template_id_5973f308___WEBPACK_IMPORTED_MODULE_0__["render"]; });\n\n/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "staticRenderFns", function() { return _node_modules_cache_loader_dist_cjs_js_cacheDirectory_node_modules_cache_vue_loader_cacheIdentifier_d12f3824_vue_loader_template_node_modules_vue_loader_lib_loaders_templateLoader_js_vue_loader_options_node_modules_cache_loader_dist_cjs_js_ref_0_0_node_modules_vue_loader_lib_index_js_vue_loader_options_ImageViewer_vue_vue_type_template_id_5973f308___WEBPACK_IMPORTED_MODULE_0__["staticRenderFns"]; });\n\n\n\n//# sourceURL=webpack:///./src/components/filebrowser/viewers/ImageViewer.vue?')}}]);