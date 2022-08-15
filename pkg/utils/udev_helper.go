/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-08-10 16:06:12
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-10 16:11:37
 * @FilePath: /CasaOS/pkg/utils/udev_helper.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package utils

// func getOptionnalMatcher() (matcher netlink.Matcher, err error) {
// 	if filePath == nil || *filePath == "" {
// 		return nil, nil
// 	}

// 	stream, err := ioutil.ReadFile(*filePath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if stream == nil {
// 		return nil, fmt.Errorf("Empty, no rules provided in \"%s\", err: %w", *filePath, err)
// 	}

// 	var rules netlink.RuleDefinitions
// 	if err := json.Unmarshal(stream, &rules); err != nil {
// 		return nil, fmt.Errorf("Wrong rule syntax, err: %w", err)
// 	}

// 	return &rules, nil
// }
